"""
This is the vsphere_virtual_disk_volume_source module in k8spkg.api.core.v1 package.
This file was generated by the KCL auto-gen tool. DO NOT EDIT.
Editing this file might prove futile when you re-run the KCL auto-gen generate command.
"""


schema VsphereVirtualDiskVolumeSource:
    """ Represents a vSphere volume resource.

    Attributes
    ----------
    fsType : str, default is Undefined, optional
         Filesystem type to mount. Must be a filesystem type supported by the host operating system. Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
    storagePolicyID : str, default is Undefined, optional
         Storage Policy Based Management (SPBM) profile ID associated with the StoragePolicyName.
    storagePolicyName : str, default is Undefined, optional
         Storage Policy Based Management (SPBM) profile name.
    volumePath : str, default is Undefined, required
         Path that identifies vSphere volume vmdk
    """


    fsType?: str

    storagePolicyID?: str

    storagePolicyName?: str

    volumePath: str


