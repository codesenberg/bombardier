-- init random
counter = 1

request = function(path, method, body)
    
    -- define the path that will search for q=%v 9%v being a random number between 0 and 1000)
    url_path = "/somepath/search?q=" .. counter
    counter = counter + 1

    -- if we want to print the path generated
    --print(url_path)

    return url_path, method, body, nil
end