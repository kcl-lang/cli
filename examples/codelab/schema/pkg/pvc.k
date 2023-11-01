import k8spkg.api.core.v1

protocol PVCProtocol:
    pvc?: {str:}

mixin PersistentVolumeClaimMixin for PVCProtocol:
    """
    PersistentVolumeClaim (PVC) sample:
    Link: https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims
    """
    
    # Mix in a new attribute `kubernetesPVC`
    kubernetesPVC?: v1.PersistentVolumeClaim

    if pvc:
        kubernetesPVC = v1.PersistentVolumeClaim {
            metadata.name = pvc.name
            metadata.labels = pvc.labels
            spec = {
                accessModes = pvc.accessModes
                resources = pvc.resources
                storageClassName = pvc.storageClassName
            }
        }
