"""
This is the flex_volume_source module in k8spkg.api.core.v1 package.
This file was generated by the KCL auto-gen tool. DO NOT EDIT.
Editing this file might prove futile when you re-run the KCL auto-gen generate command.
"""


schema FlexVolumeSource:
    """ FlexVolume represents a generic volume resource that is provisioned/attached using an exec based plugin.

    Attributes
    ----------
    driver : str, default is Undefined, required
         Driver is the name of the driver to use for this volume.
    fsType : str, default is Undefined, optional
         Filesystem type to mount. Must be a filesystem type supported by the host operating system. Ex. "ext4", "xfs", "ntfs". The default filesystem depends on FlexVolume script.
    options : {str:str}, default is Undefined, optional
         Optional: Extra command options if any.
    readOnly : bool, default is Undefined, optional
         Optional: Defaults to false (read/write). ReadOnly here will force the ReadOnly setting in VolumeMounts.
    secretRef : LocalObjectReference, default is Undefined, optional
         Optional: SecretRef is reference to the secret object containing sensitive information to pass to the plugin scripts. This may be empty if no secret object is specified. If the secret object contains more than one secret, all secrets are passed to the plugin scripts.
    """


    driver: str

    fsType?: str

    options?: {str:str}

    readOnly?: bool

    secretRef?: LocalObjectReference


