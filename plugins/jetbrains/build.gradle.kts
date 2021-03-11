import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

plugins {
    id("org.jetbrains.intellij") version "0.7.2"
    id("org.jetbrains.kotlin.jvm") version "1.4.21"
    id("org.kordamp.gradle.markdown") version "2.2.0"
    id("java")
}

group = "WhaleLint"
version = "0.0.5"

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

tasks.register<Copy>("copyChangelog") {
    from(file("$buildDir/../../vscode/CHANGELOG.md"))
    into(file("$buildDir/idea-sandbox/plugins/whalelint/docs/"))
}

tasks.markdownToHtml {
    sourceDir = file("$buildDir/idea-sandbox/plugins/whalelint/docs")
    outputDir = file("$buildDir/idea-sandbox/plugins/whalelint/docs")
}

tasks.getByName<org.jetbrains.intellij.tasks.PatchPluginXmlTask>("patchPluginXml") {
    dependsOn("copyChangelog", "markdownToHtml")

    val changelogPath = "$buildDir/idea-sandbox/plugins/whalelint/docs/CHANGELOG.html"

    if (file(changelogPath).exists()) {
        changeNotes(file(changelogPath).readText())
    }

    version("0.0.5")
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
    dependsOn("copyChangelog", "markdownToHtml")
}
