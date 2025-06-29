# Claude Code Session Prompts

This file contains all the prompts used during the Claude Code session for developing the Datadog metrics query tool.

## Session Overview
This session involved developing a Go application that queries Datadog metrics API to analyze Django request patterns using binary search algorithms.

---

## Prompt 1
```
You are a golang developer. Inspect main.go, replace the timestamp 1751061600000 with a new parameter END_DATE. Reuse the code what was implemented for START_DATE. Try to factorize generic method such as the parsing of a date. Execute the program main.go with the following environnemt variable to confirm that it works : DD_API_KEY=<some_key>;DD_APP_KEY=<some_app_key>;DD_SITE=datadoghq.com
```

## Prompt 2
```
Show the command line to execute the same script with START_DATE=2024-01-01 and END_DATE=2024-12-31
```

## Prompt 3
```
Replace "*admin*" Line 109 by a param. This parameter name is mandatory to execute the script. Replace the string, add error handler if the param is not set, and test with "*admin*" to confirm that you have the same result as before.
```

## Prompt 4
```
resourceFilter must use the function pathToDatadogMetrics to convert the param value from the command line, to a valid Datadog resource_name. Update the code.
```

## Prompt 5
```
[Request interrupted by user for tool use]Read the RESOURCE_FILTER value from the command line, convert it to an Array with pathToDatadogMetrics, then iterate over this array. For each array value, execute the call to DatadoAPI. Do not test, but output a command line to let me test.
```

## Prompt 6
```
I want to implement a dichotomic search, to find the most recent value for a datadog metric. Use start_date and end_date as a range of date. For instance if start_date is 2024-01-01 and end_date is 2025-01-01 then the range is one year. Split the range into 2 subranges. For each range, execute the call to the Datadog API. If you find a value for the metric, then repeat with a smaller range, until you find the most recent value for the metric. If you cannot find any value for this metric between start_date and end_date, then report this information as a success. If you find a value for the metric in the range, output the value and the date that was found, and report a warning.
```

## Prompt 7
```
Output today date as a timestamp. This will be the END-DATE. Substract 18 months to this END_DATE to find the START_DATE, output the value as a timestamp. This will be the START_DATE.
```

## Prompt 8
```
Update main.go. If START_DATE and END_DATE are not specified, then default the END_DATE to current date, and default START_DATE to the END_DATE minus 18 months. Update the CLAUDE.md documentation to reflect this default behavior.
```

## Prompt 9
```
Capture the result call to Datadog API and write a function to retrieve the most recent timestamp where a value is not null for each element of the "series" json array. The json result has a data object with an attributes subobject. This "attributes" object has 3 properties : series, times and  values. Series is an array that contains the name of a Datadog metric, for instance resource_name:get_bm/category/leaf/light_/. 
times is an array of timestamp, the most recent element is usually the last element of the array "times".
values is an array of arrays. If "series" had 2 datadog metric, then the values array will have 2 arrays. The first array is for the first element of "series". The second array is for the second element of "series".
I will give you a sample JSON result, I want you to write a function that accepts this JSON, and return an array with a tuple. The tuple should be the name of the group_tags and the second tuple value should be the last timestamp for which "values" is not null.
```

## Prompt 10
```
update the result json object, add a field "last_seen_date", convert the timestamp to a human readable date, use GMT as the timezone.
```

## Prompt 11
```
generate a .gitignore file for this golang project, make sure to also exclude .idea
```

## Prompt 12
```
git add *
```

## Prompt 13
```
generate a commit message that explains what does this program do
```

## Prompt 14
```
Can you print Last seen as a more human readable date? Maybe ISO-8859-15 format or French format.
```

## Prompt 15
```
can you use format DD MMMM YYYY HH:MM:SS UTC
```

## Prompt 16
```
[Prompt removed for security reason]
```

## Prompt 17
```
You are a Golang developer. I want to fix main_test.go. Can you update main_test.go to reflect the last changes from main.go?
```

## Prompt 18
```
Update main.go line 198 so that the string "service:badoom*" can be configured if needed. Add an optional param SERVICE with default value "badoom". Keep the * to match any service name that starts with the param value. Update CLAUDE.md documentation to reflect this.
```

## Prompt 19
```
Can you save all the prompts during this session to a new file claude_prompt.md?
```

---

## Summary
This session covered:
1. Parameter implementation for END_DATE and START_DATE
2. Making RESOURCE_FILTER mandatory with error handling
3. Converting URL paths to Datadog metric names
4. Implementing binary search for finding most recent metric values
5. Adding dynamic date calculation (18 months back from today)
6. JSON parsing and timestamp processing
7. Date formatting improvements
8. Git security (removing sensitive API keys from history)
9. Test file updates
10. Making service name configurable
11. Documentation updates throughout

The final result is a robust Go application that efficiently searches for the most recent activity on Datadog metrics using binary search algorithms.