---
title: Newport Reference Manual
source: therealtruex.com (recovered from web.archive.org)
---

<span class="relative-nav"> <span class="fishdown"> </span> <span class="centercomp"> </span></span>

## Newport Reference Manual

  

------------------------------------------------------------------------

<img src="static/p-identical-to-p-tiny.png" data-a="" data-by="" data-followed="" data-identical-to="" data-p.="" alt="Logo fo Newport. A p followed by a triple bar (a math symbol for " />

#### Contents

- [Time](#Time)
  - [`CURRENT-TIME`](#CURRENT-TIME)
  - [`STRING->TZ`](#STRING-TZ)
  - [`TZ->STRING`](#TZ-STRING)
  - [`+TIME-ZONES+`](#TIME-ZONES)
  - [`+WEEK-DAYS+`](#WEEK-DAYS)
  - [`+MONTH-NAMES+`](#MONTH-NAMES)
- [Structures](#Structures)
  - [`STRUCTURE-PREDICATE`](#STRUCTURE-PREDICATE)
  - [`STRUCTURE-COPIER`](#STRUCTURE-COPIER)
  - [`STRUCTURE-BOA-CONSTRUCTORS`](#STRUCTURE-BOA-CONSTRUCTORS)
  - [`STRUCTURE-KEYWORD-CONSTRUCTOR`](#STRUCTURE-KEYWORD-CONSTRUCTOR)
  - [`STRUCTURE-SLOTS`](#STRUCTURE-SLOTS)
- [Classes](#Classes)
  - [`CLASS-SLOT-INITARGS`](#CLASS-SLOT-INITARGS)
  - [`CLASS-SLOT-LIST`](#CLASS-SLOT-LIST)
  - [`CLASS-SLOT-INITARGS`](#CLASS-SLOT-INITARGS)
  - [`CLASS-SLOT-LIST`](#CLASS-SLOT-LIST)
- [System](#System)
  - [`SYSINFO`](#SYSINFO)
  - [`VARIABLE-SPECIAL-P`](#VARIABLE-SPECIAL-P)
  - [`GETENV`](#GETENV)
  - [`GETENV`](#GETENV)
  - [`DEFAULT-DIRECTORY`](#DEFAULT-DIRECTORY)
- [Pipes](#Pipes)
  - [`WITH-OPEN-PIPE`](#WITH-OPEN-PIPE)
  - [`CLOSE-PIPE`](#CLOSE-PIPE)
  - [`PIPE-INPUT`](#PIPE-INPUT)
  - [`PIPE-OUTPUT`](#PIPE-OUTPUT)
  - [`RUN-PROG`](#RUN-PROG)
- [Conditions](#Conditions)
  - [`NOT-IMPLEMENTED`](#NOT-IMPLEMENTED)
  - [`CODE`](#CODE)
- [CL Utils](#CLUtils)
  - [`ARGLIST`](#ARGLIST)
  - [`VARIABLE-NOT-SPECIAL`](#VARIABLE-NOT-SPECIAL)
  - [`MK-ARR`](#MK-ARR)
  - [`COMPOSE`](#COMPOSE)
  - [`DEFCONST`](#DEFCONST)

### <span id="Time" class="jump">Time</span>

<span id="CURRENT-TIME" class="jump">`CURRENT-TIME`</span>  
Print the current time to the stream (defaults to t).  
  
<span id="STRING-TZ" class="jump">`STRING->TZ`</span>  
Find the OBJ (symbol or string) in +TIME-ZONES+.  
  
<span id="TZ-STRING" class="jump">`TZ->STRING`</span>  
Convert the CL timezone (rational \[-24;24\], multiple of 3600) to a string.  
  
<span id="TIME-ZONES" class="jump">`+TIME-ZONES+`</span>  
The string representations of the time zones.  
  
<span id="WEEK-DAYS" class="jump">`+WEEK-DAYS+`</span>  
The names of the days of the week.  
  
<span id="MONTH-NAMES" class="jump">`+MONTH-NAMES+`</span>  
The names of the months.  
  

### <span id="Structures" class="jump">Structures</span>

<span id="STRUCTURE-PREDICATE" class="jump">`STRUCTURE-PREDICATE`</span>  
Return the structure predicate name.  
  
<span id="STRUCTURE-COPIER" class="jump">`STRUCTURE-COPIER`</span>  
Return the structure copier name.  
  
<span id="STRUCTURE-BOA-CONSTRUCTORS" class="jump">`STRUCTURE-BOA-CONSTRUCTORS`</span>  
Return the list of structure BOA constructor names.  
  
<span id="STRUCTURE-KEYWORD-CONSTRUCTOR" class="jump">`STRUCTURE-KEYWORD-CONSTRUCTOR`</span>  
Return the structure keyword constructor name.  
  
<span id="STRUCTURE-SLOTS" class="jump">`STRUCTURE-SLOTS`</span>  
Return the list of structure slot names.  
  

### <span id="Classes" class="jump">Classes</span>

<span id="CLASS-SLOT-INITARGS" class="jump">`CLASS-SLOT-INITARGS`</span>  
Return the list of initargs of a CLASS. CLASS can be a symbol, a class object (as returned by \`class-of') or an instance of a class. If the second optional argument ALL is non-NIL (default), initargs for all slots are returned, otherwise only the slots with :allocation type :instance are returned.  
  
<span id="CLASS-SLOT-LIST" class="jump">`CLASS-SLOT-LIST`</span>  
Return the list of slots of a CLASS. CLASS can be a symbol, a class object (as returned by \`class-of') or an instance of a class. If the second optional argument ALL is non-NIL (default), all slots are returned, otherwise only the slots with :allocation type :instance are returned.  
  

### <span id="System" class="jump">System</span>

<span id="SYSINFO" class="jump">`SYSINFO`</span>  
Print the current environment to a stream.  
  
  
<span id="VARIABLE-SPECIAL-P" class="jump">`VARIABLE-SPECIAL-P`</span>  
Return T if the symbol names a global special variable.  
  
<span id="GETENV" class="jump">`GETENV`</span>  
Set an environment variable.  
  
<span id="GETENV" class="jump">`GETENV`</span>  
Return the value of the environment variable.  
  
<span id="DEFAULT-DIRECTORY" class="jump">`DEFAULT-DIRECTORY`</span>  
The default directory.  
  

### <span id="Pipes" class="jump">Pipes</span>

<span id="WITH-OPEN-PIPE" class="jump">`WITH-OPEN-PIPE`</span>  
Open the pipe, do something, then close it.  
  
<span id="CLOSE-PIPE" class="jump">`CLOSE-PIPE`</span>  
Close the pipe stream.  
  
<span id="PIPE-INPUT" class="jump">`PIPE-INPUT`</span>  
Return an input stream from which the command output will be read.  
  
<span id="PIPE-OUTPUT" class="jump">`PIPE-OUTPUT`</span>  
Return an output stream which will go to the command.  
  
<span id="RUN-PROG" class="jump">`RUN-PROG`</span>  
Common interface to shell. Does not return anything useful.  
  

### <span id="Conditions" class="jump">Conditions</span>

<span id="NOT-IMPLEMENTED" class="jump">`NOT-IMPLEMENTED`</span>  
Your implementation does not support this functionality.  
  
<span id="CODE" class="jump">`CODE`</span>  
An error in the user code.  
  

### <span id="CLUtils" class="jump">CL Utils</span>

<span id="ARGLIST" class="jump">`ARGLIST`</span>  
Return the signature of the function.  
  
<span id="VARIABLE-NOT-SPECIAL" class="jump">`VARIABLE-NOT-SPECIAL`</span>  
Undo the global special declaration. This returns a \_new\_ symbol with the same name, package, fdefinition, and plist as the argument. This can be confused by imported symbols. Also, (FUNCTION-LAMBDA-EXPRESSION (FDEFINITION NEW)) will return the OLD (uninterned!) symbol as its 3rd value. BEWARE!  
<span id="MK-ARR" class="jump">`MK-ARR`</span>  
Make array with elements of TYPE, initializing.  
  
<span id="COMPOSE" class="jump">`COMPOSE`</span>  
Make a new function by composition of previous functions.  
  
<span id="DEFCONST" class="jump">`DEFCONST`</span>  
Define a typed constant.

## Acknowledgements

To the CL Gardeners project (now defunct) which suggested this activity.  
To:  
bordeaux-threads (threading)  
usocket (sockets)  
[CLOCC's obsolete PORT](http://clocc.sourceforge.net/dist/port.html)  
[trivial gray streams](https://github.com/trivial-gray-streams/trivial-gray-streams)  
trivial garbage  
closer-mop  
[cl-fad (pathnames)](https://github.com/edicl/cl-fad)  
for obsoleting \*most\* of CLOCC PORT. <span class="fishup"> </span>
