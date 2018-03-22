# Scripting 

You can use lua scripting to be able to change the path, method and body of your request.
It can be useful if you are testing a path that do cache. Or if you want to benchmark many path of your site, to know the benchmark of your site.

You can find here 3 scrits already prepare to tests :

- incremental-path.lua that increments a number in the path at each request
- random-path.lua that generates a random number at each request (to test don't forget to change he path for the paths.txt file)
- load-file-path.lua that load a file of path


```bash
bombardier -n 10 -l -z /path-to-scripts/my-script.lua http://localhost/
```

If you want to check what are the path generated you can decomment the print(url_path) on the scripts, it will print in the stdout.
