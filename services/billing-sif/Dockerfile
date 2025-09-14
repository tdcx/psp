# ---- build ----
FROM eclipse-temurin:21-jdk-jammy AS build
WORKDIR /workspace
# Copia gradle wrapper y config primero para cache
COPY services/billing-sif/gradlew services/billing-sif/gradle/ ./ 
COPY services/billing-sif/build.gradle* services/billing-sif/settings.gradle* ./ 
# Dependencias (si usas versiones catalog, copia gradle/libs.versions.toml)
COPY services/billing-sif/gradle/libs.versions.toml gradle/libs.versions.toml
# Descarga dependencias
RUN ./gradlew --no-daemon dependencies || true
# Copia el resto
COPY services/billing-sif/ .
# Construye fat jar (ajusta tarea si usas Shadow o Spring Boot)
# Para Spring Boot: ./gradlew bootJar
RUN ./gradlew clean build --no-daemon
# Detecta jar resultante
RUN mkdir -p /out && \
    cp $(find build/libs -maxdepth 1 -name "*.jar" -type f | head -n 1) /out/app.jar

# ---- runtime ----
# Distroless Java 21, usuario no root por defecto
FROM gcr.io/distroless/java21-debian12:nonroot
WORKDIR /app
COPY --from=build /out/app.jar /app/app.jar
EXPOSE 8080
# Si tu app expone /health, ajusta puerto/path en chart/values
ENTRYPOINT ["java","-jar","/app/app.jar"]
