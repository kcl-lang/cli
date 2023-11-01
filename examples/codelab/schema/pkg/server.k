schema Server(Deployment):
    mixin [PersistentVolumeClaimMixin]
    pvc?: {str:}
    """ pvc user interface data defined by PersistentVolumeClaimMixin """
