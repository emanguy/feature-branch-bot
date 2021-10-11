FROM gradle:7.2-jdk11 AS build-container
COPY src ./src
COPY ["build.gradle.kts", "gradle.properties", "settings.gradle.kts", "./"]
RUN gradle shadowJar

FROM openjdk:11-jre-slim-bullseye
RUN mkdir -p /opt/feature-branch-bot
WORKDIR /opt/feature-branch-bot/
COPY ./bot-startup.sh ./
RUN chmod u+x ./bot-startup.sh
COPY --from=build-container /home/gradle/build/libs/FeatureBranchBot.jar ./
CMD ./bot-startup.sh
