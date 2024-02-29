
# Add e2e test

## 1. Quick Start
We have prepared the `e2e-init` command to help you initialize an empty e2e test case. The command format is as follows:

```
make e2e-init TS=<test suite name>
``` 

`<test suite name>` is the name of the test case.

For example, to add a test case named `test_kcl_mod_init_1` to test the output of the `kcl mod init` command in a KCL package directory, you can use the following command:

### 1.1. Create an empty test case
```
make e2e-init TS=test_kcl_mod_init
```

You will get the following result:

```
Test suite created successfully in ./test/e2e/test_suites/test_kcl_mod_init_1.
```
By this command, we can see that a new test case `test_kcl_mod_init_1` has been created in the `test/e2e/test_suites` directory of the project.

```
$ tree test/e2e/test_suites/test_kcl_mod_init_1
test_kcl_mod_init_1
├── input # The input of the test case
├── stderr # The expected output in stderr
├── stdout # The expected output in stdout
└── test_space # The test space is a empty directory and the `input` will run in this directory
```

### 1.2. Add the test environment to the `test_space` directory.

For `test_kcl_mod_init_1`, we need to create a KCL package directory in the `test_space` directory. The directory structure is as follows:

```
$tree test/e2e/test_suites/test_kcl_mod_init_1/test_space
test_space
├── kcl.mod
├── kcl.mod.lock
└── main.k
```

### 1.3. Add the command to be tested to the `input` file. 

The content of the `input` file is as follows:

```
kcl mod init
```

### 1.4. The expected output 

`stderr` is empty, and the content in `stdout` is as follows:

```
creating new :<workspace>/test_space/kcl.mod
'<workspace>/test_space/kcl.mod' already exists
creating new :<workspace>/test_space/kcl.mod.lock
'<workspace>/test_space/kcl.mod.lock' already exists
creating new :<workspace>/main.k
'<workspace>/main.k' already exists
package 'test_space' init finished
```

In the content of the `stdout` file, `<workspace>` is a variable that will be replaced with the absolute path of the `test_space` in the test process.

### 1.5. Start test

After completing these steps, run `make e2e` to start the test process.

The test results will be displayed as follows:

```
Ran 3 of 3 Specs in 0.102 seconds
SUCCESS! -- 3 Passed | 0 Failed | 0 Pending | 0 Skipped
```

## 2. Adjust the execution path of the tested command

By default, the command in the `input` file is executed in the `test_space` directory. However, in some more complex cases, we may want to execute the command in a subdirectory of `test_space`. We can use the `conf.json` file to specify the directory in which the command is executed.

Take `test_kcl_mod_add_local` as an example to show the process. In this test case, we want to test the process of adding a KCL package in the local directory as a third-party dependency through the `kcl mod add` command.

### 2.1. Use `e2e-init` to create a test case

```
make e2e-init TS=test_kcl_mod_add_local
```

### 2.2. Prepare the test environment

In the `test_space` directory, prepare two KCL packages. The `pkg` package will add the `dep` package as a dependency through the `kcl mod add` command.

```
test_space
├── dep
 │   ├── kcl.mod
 │   ├── kcl.mod.lock
 │   └── main.k
└── pkg
    ├── kcl.mod
    ├── kcl.mod.lock
    └── main.k
```

### 2.3. Configure the execution directory

Add a `conf.json` file to the `test_kcl_mod_add_local` directory to specify the execution directory of the command.

```
{
    "cwd": “test_space/pkg"
}
```

### 2.4. Add the command to be tested to the `input` file

The content of the `input` file is as follows:

```
kcl mod add <workspace>/dep
```

### 2.5. Add the expected output

Add the expected output to the `stdout` file.

```
 adding dependency 'dep'
 add dependency 'dep:0.0.1' successfully
``` 