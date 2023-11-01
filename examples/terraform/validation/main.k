schema TFPlan:
    # Omit other attributes
    [...str]: any
    resource_changes?: [AcceptableChange]

schema AcceptableChange:
    # Omit other attributes
    [...str]: any
    check:
        # Reject AWS autoscaling group Resource delete action
        all action in change.actions {
            action not in ["delete"]
        } if type == "aws_autoscaling_group", "Disable AWS autoscaling group resource delete action for the resource ${type} ${name}"
