import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

plugins {
    kotlin("jvm") version "1.5.31"
    kotlin("plugin.serialization") version "1.5.31"
    application
}

group = "me.erittenhouse"
version = "0.0.1"

repositories {
    mavenCentral()
    maven("https://repo.eclipse.org/content/groups/releases/")
}

dependencies {
    implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:1.3.0")
    implementation("org.gitlab4j:gitlab4j-api:4.18.0")
    implementation("org.eclipse.jgit:org.eclipse.jgit:5.13.0.202109080827-r")
    implementation("org.eclipse.jgit:org.eclipse.jgit.ssh.jsch:5.13.0.202109080827-r")
    implementation("com.jcraft:jsch:0.1.55")
    testImplementation(kotlin("test"))
}

tasks.test {
    useJUnitPlatform()
}

kotlin {
    sourceSets.all {
        languageSettings.apply {
            languageVersion = "1.6"
        }
    }
}

tasks.withType<KotlinCompile>() {
    kotlinOptions.jvmTarget = "11"
}

application {
    mainClass.set("me.erittenhouse.featurebranchbot.MainKt")
}