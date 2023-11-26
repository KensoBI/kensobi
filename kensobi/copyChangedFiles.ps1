robocopy "public" "..\public" /s
robocopy "pkg" "..\pkg" /s
robocopy "packages" "..\packages" /s
robocopy "." ".." "Dockerfile"
robocopy "." ".." "README.md"