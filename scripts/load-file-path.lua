-- Load URL paths from the file
function load_url_paths_from_file(file)
    lines = {}
  
    -- Check if the file exists
    -- Resource: http://stackoverflow.com/a/4991602/325852
    local f=io.open(file,"r")
    if f~=nil then
      io.close(f)
    else
      -- Return the empty array
      return lines
    end
  
    -- If the file exists loop through all its lines
    -- and add them into the lines array
    for line in io.lines(file) do
      if not (line == '') then
        lines[#lines + 1] = line
      end
    end
  
    return lines
  end
  
  -- Load URL paths from file
  -- CHANGE BY YOUR PATH
  paths = load_url_paths_from_file("./scripts/paths.txt")
  
  print("multiplepaths: Found " .. #paths .. " paths")
  
  -- Initialize the paths array iterator
  counter = 0
  
  request = function(path, method, body)
    -- Get the next paths array element
    url_path = paths[counter]
    if url_path == nil or url_path == "" then
        url_path = "/"
    end
    counter = counter + 1
  
    -- If the counter is longer than the paths array length then reset it
    if counter > #paths then
      counter = 0
    end

    if url_path == nil or url_path == "" then
        url_path = "/"
    end

    --print(url_path)
  
    -- Return the request object with the current URL path
    return url_path, method, body, nil
  end