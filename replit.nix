{ pkgs }: {
    deps = [
        pkgs.ffmpeg.bin
        pkgs.go
        pkgs.gopls
    ];
}