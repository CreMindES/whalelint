import org.jetbrains.kotlin.gradle.tasks.KotlinCompile
import org.jetbrains.kotlin.platform.TargetPlatformVersion.NoVersion.description

plugins {
    id("org.jetbrains.intellij") version "0.7.2"
    id("org.jetbrains.kotlin.jvm") version "1.4.21"
    id("org.kordamp.gradle.markdown") version "2.2.0"
    id("java")
}

group = "WhaleLint"
version = "0.0.6"

description = "WhaleLint is a Dockerfile linter written in Golang."

repositories {
    mavenCentral()
}

dependencies {
    implementation(kotlin("stdlib-jdk8"))
    implementation("com.google.code.gson:gson:2.8.6" )
    implementation("org.jetbrains:annotations:20.1.0")
    runtimeOnly(group = "commons-io", name = "commons-io", version = "2.6")
}

// See https://github.com/JetBrains/gradle-intellij-plugin/
intellij {
    version = "2020.3"
    pluginName = "whalelint"
    updateSinceUntilBuild = false
    // setPlugins("Docker:$version")
}

tasks.register<Copy>("copyChangelogAndReadme") {
    from(file("$buildDir/../../vscode/CHANGELOG.md"), file("readme.md"))
    into(file("$buildDir/idea-sandbox/plugins/whalelint/docs/"))
}

tasks.markdownToHtml {
    sourceDir = file("$buildDir/idea-sandbox/plugins/whalelint/docs")
    outputDir = file("$buildDir/idea-sandbox/plugins/whalelint/docs")
}

tasks.getByName<org.jetbrains.intellij.tasks.PatchPluginXmlTask>("patchPluginXml") {
    dependsOn("copyChangelogAndReadme", "markdownToHtml")

    val changelogPath = "$buildDir/idea-sandbox/plugins/whalelint/docs/CHANGELOG.html"
    val readmePath    = "$buildDir/idea-sandbox/plugins/whalelint/docs/readme.html"


    if (file(changelogPath).exists()) {
        changeNotes(file(changelogPath).readText())
    }
    if (file(readmePath).exists()) {
        pluginDescription(file(readmePath).readText().replace(
            "<h1>WhaleLint JetBrains Plugin</h1>", "").replace(
            "<h2>Introduction</h2>", ""))
    }

    version("0.0.6")
}

tasks.withType<JavaCompile> {
    sourceCompatibility = "1.8"
    targetCompatibility = "1.8"
}

listOf("compileKotlin", "compileTestKotlin").forEach {
    tasks.getByName<KotlinCompile>(it) {
        kotlinOptions.jvmTarget = "1.8"
    }
}

tasks.register("copyWhaleLintBinary") {
    doLast {
        copy {
            from("$projectDir/../../../whalelint/whalelint")
            into("$buildDir/idea-sandbox/plugins/whalelint/bin/")
        }
    }
}

tasks.named("prepareSandbox") {
    finalizedBy("copyWhaleLintBinary")
}

tasks.buildPlugin {
    dependsOn("copyChangelogAndReadme", "markdownToHtml")
}
