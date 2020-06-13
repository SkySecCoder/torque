package main

import (
    "os"
    "strings"
    "torque/authMFA"
    "torque/customTypes"
    "torque/keyRotation"
    "torque/programHelp"
    "torque/assumeRole"
)

type CredDict = customTypes.CredDict

func main() {
    progArgs := os.Args
    firstArgument := strings.Split(progArgs[0], "/")
    programName := firstArgument[len(firstArgument)-1]

    if len(os.Args) == 1 || os.Args[1] == "help" || os.Args[1] == "-h" {
        programHelp.ProgramHelp(programName)
    } else if os.Args[1] == "auth" && len(os.Args) == 3 {
        authMFA.AuthMFA(progArgs[2], "")
    } else if os.Args[1] == "rotate" && len(os.Args) == 3 {
        keyRotation.RotateKey(progArgs[2])
    } else if os.Args[1] == "assume" && len(os.Args) == 3 {
        assumeRole.AssumeRole(progArgs[2])
    } else {
    programHelp.ProgramHelp(programName)
    }
}
