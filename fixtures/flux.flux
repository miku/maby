default in = "a.mab";

in | open-file | as-records | decode-mab | encode-json | write("stdout");
