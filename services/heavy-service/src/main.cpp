#include <cstdlib>
#include <iostream>
#include <string>

int main() {
    const char* envPort = std::getenv("HEAVY_SERVICE_PORT");
    std::string port = envPort ? envPort : "8083";

    std::cout << "heavy-service bootstrap started on port " << port << "\n";
    std::cout << "Next step: wire Drogon app and expose /health endpoint." << "\n";
    return 0;
}
