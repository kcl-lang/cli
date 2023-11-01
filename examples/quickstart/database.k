schema DatabaseConfig:
    enabled: bool = True
    ports: [int] = [8000, 8001, 8002]
    data: [[str|float]] = [["delta", "phi"], [3.14]]
    temp_targets: {str: float} = {cpu = 79.5, case = 72.0}

database = DatabaseConfig {
    ports = [2020, 2021]
}
