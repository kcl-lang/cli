import .my_config

# debugging
print(my_config.labels) # debugging by print

# test
assert len(my_config.labels) > 0, "labels can't be empty" # use len() to get list length
assert "env" in my_config.labels, "env label is a must"
assert my_config.cpu >= 256, "cpu cannot be less than256"
