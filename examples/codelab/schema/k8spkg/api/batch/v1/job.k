"""
This is the job module in k8spkg.api.batch.v1 package.
This file was generated by the KCL auto-gen tool. DO NOT EDIT.
Editing this file might prove futile when you re-run the KCL auto-gen generate command.
"""
import k8spkg.apimachinery.apis


schema Job:
    """ Job represents the configuration of a single job.

    Attributes
    ----------
    apiVersion : "batch/v1", default is "batch/v1", required
         APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
    kind : "Job", default is "Job", required
         Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
    metadata : apis.ObjectMeta, default is Undefined, optional
         Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
    spec : JobSpec, default is Undefined, optional
         Specification of the desired behavior of a job. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
    """


    apiVersion: "batch/v1" = "batch/v1"

    kind: "Job" = "Job"

    metadata?: apis.ObjectMeta

    spec?: JobSpec


