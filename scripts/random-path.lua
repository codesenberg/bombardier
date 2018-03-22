-- init random
math.randomseed(os.time())

request = function(path, method, body)
    
    -- define the path that will search for q=%v 9%v being a random number between 0 and 1000)
    url_path = "/somepath/search?q=" .. math.random(0,1000)

    -- if we want to print the path generated
    --print(url_path)

    return url_path, method, body, nil
end