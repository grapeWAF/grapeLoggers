# -*- coding: utf-8 -*- 
import os
import os.path
import shutil
import time
import datetime
import string
import sys
import platform
import time
import signal
import imp

imp.reload(sys)

source_dir = os.getcwd()
go_dir = r"go"

is_public_Rel = True
ISOTIMEFORMAT = '%Y-%m-%d %X'

buildfiles = [
    # srcdir,buildname isRel isLinux is64Bit
    #[r"/guard","guard_Dbg",False,True,True],
    [r"/guard","guard",True,True,True],
    #[r"/guardCertZ","guardcert_Dbg",False,True,True],
    #[r"/guardCertZ","guardcert",True,True,True],
    #[r"/guardCtrl","guardctrl_Dbg",False,True,True],
    #[r"/guardCtrl","guardctrl",True,True,True],
]

def logtofile(str):
    type = sys.getfilesystemencoding()
    timeStr = time.strftime(ISOTIMEFORMAT, time.localtime(time.time()))
    print("[" + timeStr + "] " + str)#.decode('utf-8').encode(type))
    #os.system("echo " + "["+ timeStr + "] " + str.decode('utf-8').encode(type) + " >> " + getLogPath())


def copyFiles(sourceDir,  targetDir): 
    if sourceDir.find(".svn") > 0: 
        return 
    for file in os.listdir(sourceDir): 
        sourceFile = os.path.join(sourceDir,  file) 
        targetFile = os.path.join(targetDir,  file) 
        if os.path.isfile(sourceFile): 
            if not os.path.exists(targetDir):  
                os.makedirs(targetDir)  
            if not os.path.exists(targetFile) or(os.path.exists(targetFile) and (os.path.getsize(targetFile) != os.path.getsize(sourceFile))):  
                    open(targetFile, "wb").write(open(sourceFile, "rb").read()) 
        if os.path.isdir(sourceFile): 
            First_Directory = False 
            copyFiles(sourceFile, targetFile)


def moveFileto(sourceDir,  targetDir): 
    shutil.copy(sourceDir,  targetDir)

def build_go(sourceDir,outName,is64Bit,isRel,isLinux):
    if is64Bit:
        os.environ["GOARCH"] = "amd64"
    else:
        os.environ["GOARCH"] = "386"

    if isLinux:
        os.environ["CGO_ENABLED"] = "0"
        os.environ["GOEXE"] = ""
        os.environ["GOOS"] = "linux"
    else:
        os.environ["CGO_ENABLED"] = "1"
        os.environ["GOEXE"] = ".exe"
        os.environ["GOOS"] = "windows"

    os.environ["CC"] = "gcc"
    os.environ["CXX"] = "g++"

    ccurdir = os.curdir
    os.chdir(sourceDir)

    goRunBuild = go_dir
    goRunBuild += " build -v"

    if isRel:
        goRunBuild += ' -ldflags \"-w -s\"'

    goRunBuild += " -o " + outName
    os.system(goRunBuild)
    os.chdir(ccurdir)
    logtofile(goRunBuild)
    
def build_all():
    for bfile in buildfiles:
        build_go(source_dir + bfile[0],bfile[1],bfile[4],bfile[2],bfile[3])

if  __name__ =="__main__": 
    logtofile("star build all...")
    build_all()
    logtofile("release done!!!")