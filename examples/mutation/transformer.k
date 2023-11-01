items = lambda {
    # Construct resource and params
    resource = option("resource_list")
    items = resource.items
    params = resource.functionConfig.spec.params
    new_replicas: int = params.replicas
    min_replicas: int = params.min_replicas or 0
    max_replicas: int = params.max_replicas or 99999
    # Define the validation function
    validate_replica_limit = lambda item, min_replicas: int, max_replicas: int, new_replicas: int {
        replicas = item.spec.replicas or 0
        assert min_replicas < replicas < max_replicas, "The provided number of replicas ${replicas} is not allowed for ${item.kind}: ${item.metadata.name}. Allowed range: ${min} - ${max}"
        item | {
            if typeof(new_replicas) == "int":
                spec.replicas = new_replicas
        }
    }
    # Validate All resource
    [validate_replica_limit(i, min_replicas, max_replicas, new_replicas) for i in items]
}()
