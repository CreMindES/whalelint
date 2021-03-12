package whalelint

import com.intellij.psi.PsiFile
import com.intellij.psi.PsiDocumentManager
import com.intellij.codeInsight.intention.IntentionAction
import com.intellij.codeInspection.util.IntentionName
import com.intellij.codeInspection.util.IntentionFamilyName
import com.intellij.ide.BrowserUtil.browse
import com.intellij.ide.plugins.PluginManagerCore
import com.intellij.lang.annotation.*
import com.intellij.notification.NotificationGroupManager
import com.intellij.notification.NotificationType
import com.intellij.openapi.application.PathManager
import com.intellij.openapi.application.runReadAction
import com.intellij.openapi.diagnostic.Logger
import com.intellij.openapi.util.TextRange
import com.intellij.openapi.editor.Editor
import com.intellij.openapi.extensions.PluginId
import com.intellij.openapi.project.Project
import com.intellij.openapi.util.io.FileUtil
import com.intellij.util.IncorrectOperationException
import org.apache.commons.io.IOUtils
import kotlin.Throws

import java.io.BufferedReader
import java.io.File
import java.io.InputStreamReader
import java.nio.file.Files
import java.nio.file.Paths

class WhaleLintExternalAnnotator : ExternalAnnotator<PsiFile, List<ValidationResult.Issue>>() {
    override fun getPairedBatchInspectionShortName(): String {
        return Bundle.message("inspection.short.name")
    }

    override fun collectInformation(psiFile: PsiFile, editor: Editor, hasErrors: Boolean): PsiFile {
        return psiFile
    }

    override fun doAnnotate(psiFile: PsiFile): List<ValidationResult.Issue>? {

        val text = runReadAction { if (psiFile.isValid) psiFile.text else null } ?: return null
        val file = psiFile.virtualFile ?: return null
        val copy = createTempFile(text.toByteArray(file.charset))

        // get path to WhaleLint executable bundled with the plugin
        val pluginPath = PluginManagerCore.getPlugin(PluginId.getId("tamas_g_barna.whalelint"))!!.pluginPath
        val executable = Paths.get(pluginPath.toString(), "bin/whalelint")

        // assemble WhaleLint calling command
        val command = arrayOf(executable.toString(), copy.absolutePath, "--format=json")

        // Workaround for silly-silly JetBrains system bug.
        // For details, please see https://github.com/intellij-rust/intellij-rust/pull/6869, or
        // https://youtrack.jetbrains.com/issue/TW-4651, or any other similar ticket on the subject ...
        if (!Files.isExecutable(executable)) {
            executable.toFile().setExecutable(true)
        }

        try {
            LOG.info("WhaleLint | Running ...")

            // execute WhaleLint
            val process = Runtime.getRuntime().exec(command)
            val reader = BufferedReader(InputStreamReader(process.inputStream))

            // read result of linting
            var line: String? = ""
            val output = StringBuffer(110)
            while (reader.readLine().also { line = it } != null) {
                output.append(line)
            }
            val findingsStr = output.toString()

            val stdErr = IOUtils.toString(process.errorStream, Charsets.UTF_8)

            // basic result checks
            if (findingsStr.isEmpty()) {
                return null
            } else if (stdErr.isNotEmpty()) {
                notifyError(psiFile.project, "Error: $stdErr", NotificationType.ERROR)
            }

            // convert JSON str to GSON object
            return ValidationResult.parse(findingsStr)
        } catch (e: Exception) {
            LOG.error("Error running inspection: ", e)
            notifyError(psiFile.project, "Error running inspection: ${e.message}", NotificationType.ERROR)
        } finally {
            FileUtil.delete(copy)
        }

        return null
    }

    override fun apply(file: PsiFile, validationResult: List<ValidationResult.Issue>?, holder: AnnotationHolder) {
        // nothing to do in case there is no lint finding
        if (validationResult == null) {
            return
        }

        val document = PsiDocumentManager.getInstance(file.project).getDocument(file) ?: return
        for (issue in validationResult) {
            val severity = getHighlightSeverity(issue)

            if (issue.locationRange!!.start!!.lineNumber == 0) {
                val a = issue.locationRange!!.start.toString()
                println(a)
            }
            val lineStartOffset = document.getLineStartOffset(issue.locationRange!!.start!!.lineNumber - 1) +
                    issue.locationRange!!.start!!.charNumber
            val highlightLength = issue.locationRange!!.end!!.charNumber -
                    issue.locationRange!!.start!!.charNumber
            val lineEndOffset = lineStartOffset + highlightLength

            val myFix: IntentionAction = object : IntentionAction {
                override fun getText(): @IntentionName String {
                    return "Open Docs in browser."
                }

                override fun getFamilyName(): @IntentionFamilyName String {
                    return "FamilyName"
                }

                override fun isAvailable(project: Project, editor: Editor, file: PsiFile): Boolean {
                    return true
                }

                @Throws(IncorrectOperationException::class)
                override fun invoke(project: Project, editor: Editor, file: PsiFile) {
                    browse("https://github.com/CreMindES/whalelint/tree/main/docs/rule/set/" +
                            "${issue.rule!!.ruleId!!.toLowerCase()}.md")
                }

                override fun startInWriteAction(): Boolean {
                    return false
                }
            }

            val textRange = TextRange.create(lineStartOffset, lineEndOffset)
            val message = issue.message + " (WL " + issue.rule!!.ruleId + ")"
            val builder = holder.newAnnotation(severity, message)
                    .range(textRange)
                    .problemGroup { "TestProblemName" }
                    .withFix(myFix)
            builder.create()
        }
    }

    companion object {
        private val LOG = Logger.getInstance(WhaleLintExternalAnnotator::class.java)

        private fun getHighlightSeverity(issue: ValidationResult.Issue): HighlightSeverity {
            return when (issue.rule!!.severity!!.toLowerCase()) {
                "deprecation" -> HighlightSeverity.WEAK_WARNING
                "error" -> HighlightSeverity.ERROR
                "info" -> HighlightSeverity.INFORMATION
                "warning" -> HighlightSeverity.WARNING
                else -> HighlightSeverity.INFORMATION
            }
        }
    }

    private fun createTempFile(bytes: ByteArray): File {
        val tempFile = FileUtil.createTempFile(File(PathManager.getTempPath()), "tmp_whalelint", "")
        tempFile.writeBytes(bytes)
        return tempFile
    }

    private fun notifyError(project: Project, content: String, notificationType: NotificationType) {
        NotificationGroupManager.getInstance().getNotificationGroup("WhaleLint Notification Group")
            .createNotification(content, notificationType)
            .notify(project)
    }
}
