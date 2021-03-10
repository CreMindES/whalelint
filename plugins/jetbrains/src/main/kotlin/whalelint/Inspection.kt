package whalelint

import com.intellij.codeInspection.LocalInspectionTool
import com.intellij.codeInspection.ex.ExternalAnnotatorBatchInspection
import com.intellij.psi.PsiFile
import com.intellij.codeInspection.InspectionManager
import com.intellij.codeInspection.ProblemDescriptor
import com.intellij.codeInspection.ExternalAnnotatorInspectionVisitor
import com.intellij.codeInspection.ProblemsHolder
import com.intellij.codeInspection.LocalInspectionToolSession
import com.intellij.openapi.diagnostic.Logger
import com.intellij.psi.PsiElementVisitor
import com.intellij.psi.PsiElement

class Inspection : LocalInspectionTool(), ExternalAnnotatorBatchInspection {
    override fun getDisplayName(): String {
        return Bundle.message("inspection.display.name")
    }

    override fun getShortName(): String {
        return Bundle.message("inspection.short.name")
    }

    override fun isEnabledByDefault(): Boolean {
        return true
    }

    override fun checkFile(psiFile: PsiFile, manager: InspectionManager, isOnTheFly: Boolean): Array<ProblemDescriptor> {
        LOG.debug("Inspection has been invoked for " + psiFile.name)
        return ExternalAnnotatorInspectionVisitor.checkFileWithExternalAnnotator(psiFile, manager, isOnTheFly, WhaleLintExternalAnnotator())
    }

    override fun buildVisitor(holder: ProblemsHolder, isOnTheFly: Boolean, session: LocalInspectionToolSession): PsiElementVisitor {
        val externalAnnotator = WhaleLintExternalAnnotator()
        return ExternalAnnotatorInspectionVisitor(holder, externalAnnotator, isOnTheFly)
    }

    override fun isSuppressedFor(element: PsiElement): Boolean {
        return false
    }

    companion object {
        private val LOG = Logger.getInstance(Inspection::class.java)
    }
}
