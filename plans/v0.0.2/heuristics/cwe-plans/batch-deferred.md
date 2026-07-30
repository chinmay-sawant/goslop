# Batch deferred — Catalogue Deferred

> **Parent:** `plans/v0.0.2/heuristics/cwe-plans/README.md` — v0.0.2 Python CWE remaining heuristics
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#52](https://github.com/chinmay-sawant/goslop/issues/52) expansion
> **Status:** **deferred** — no implement rows; catalogue honesty only
> **IDs (185):** CWE-115, CWE-134, CWE-166, CWE-167, CWE-168, CWE-183, CWE-190, CWE-193, CWE-205, CWE-222, CWE-223, CWE-224, CWE-229, CWE-233, CWE-237, CWE-241, CWE-253, CWE-262, CWE-263, CWE-266, CWE-267, CWE-268, CWE-270, CWE-273, CWE-274, CWE-277, CWE-278, CWE-280, CWE-281, CWE-283, CWE-294, CWE-301, CWE-303, CWE-305, CWE-308, CWE-309, CWE-322, CWE-324, CWE-325, CWE-341…
> **PR policy:** one PR for this batch only — do not mix other wave IDs

---

## Architecture constraints

| Constraint | Detail |
|------------|--------|
| Package | `internal/lang/python/detectors/cwe/` only |
| Registration | `RegisterRule("CWE-*", detect…, &Meta…, gates…)` from `init()` |
| Scan | Existing `PyCweScan` — no new detector type |
| Detection | Pure-Go source patterns + needles / `PyCweFacts` |
| Language | `LanguagePython` gate already in scan |
| Plugin | **Do NOT invent a second plugin.** `detectors.All()` already wires CWE scan |
| IDs | Always `CWE-*`; metadata aligned with chunk catalogue |
| **File size policy** | Target max **1500** / hard max **2000** per Go file |
| Target file(s) | `(none — no RegisterRule)` |
| Validation | `make lint` + `make test` unchecked until green |
| Fixtures | `tests/fixtures/python/cwe/CWE-N-{{vulnerable,safe}}.txt` |

## Overview

Design/process/architectural or no clear pure-Go source sink in v0

| Rule | Name | Tier | Category | Relevance | Chunk |
|------|------|------|----------|-----------|-------|
| `CWE-115` | Misinterpretation of Input | C | General | Medium | `cwe-101-150.json` |
| `CWE-134` | Use of Externally-Controlled Format String | C | Logging | Medium | `cwe-101-150.json` |
| `CWE-166` | Improper Handling of Missing Special Element | C | General | Medium | `cwe-151-200.json` |
| `CWE-167` | Improper Handling of Additional Special Element | C | General | Medium | `cwe-151-200.json` |
| `CWE-168` | Improper Handling of Inconsistent Special Elements | C | General | Medium | `cwe-151-200.json` |
| `CWE-183` | Permissive List of Allowed Inputs | C | General | Medium | `cwe-151-200.json` |
| `CWE-190` | Integer Overflow or Wraparound | C | Number Processing | Medium | `cwe-151-200.json` |
| `CWE-193` | Off-by-one Error | C | Error Handling | Medium | `cwe-151-200.json` |
| `CWE-205` | Observable Behavioral Discrepancy | C | General | Medium | `cwe-201-250.json` |
| `CWE-222` | Truncation of Security-relevant Information | C | General | Medium | `cwe-201-250.json` |
| `CWE-223` | Omission of Security-relevant Information | C | General | Medium | `cwe-201-250.json` |
| `CWE-224` | Obscured Security-relevant Information by Alternate Name | C | General | Medium | `cwe-201-250.json` |
| `CWE-229` | Improper Handling of Values | C | General | Medium | `cwe-201-250.json` |
| `CWE-233` | Improper Handling of Parameters | C | General | Medium | `cwe-201-250.json` |
| `CWE-237` | Improper Handling of Structural Elements | C | General | Medium | `cwe-201-250.json` |
| `CWE-241` | Improper Handling of Unexpected Data Type | C | General | Medium | `cwe-201-250.json` |
| `CWE-253` | Incorrect Check of Function Return Value | C | General | Medium | `cwe-251-300.json` |
| `CWE-262` | Not Using Password Aging | C | Security | Medium | `cwe-251-300.json` |
| `CWE-263` | Password Aging with Long Expiration | C | Security | Medium | `cwe-251-300.json` |
| `CWE-266` | Incorrect Privilege Assignment | C | Security | Medium | `cwe-251-300.json` |
| `CWE-267` | Privilege Defined With Unsafe Actions | C | Security | Medium | `cwe-251-300.json` |
| `CWE-268` | Privilege Chaining | C | Security | Medium | `cwe-251-300.json` |
| `CWE-270` | Privilege Context Switching Error | C | Security | Medium | `cwe-251-300.json` |
| `CWE-273` | Improper Check for Dropped Privileges | C | Security | Medium | `cwe-251-300.json` |
| `CWE-274` | Improper Handling of Insufficient Privileges | C | Security | Medium | `cwe-251-300.json` |
| `CWE-277` | Insecure Inherited Permissions | C | Security | Medium | `cwe-251-300.json` |
| `CWE-278` | Insecure Preserved Inherited Permissions | C | Security | Medium | `cwe-251-300.json` |
| `CWE-280` | Improper Handling of Insufficient Permissions or Privileges | C | Security | Medium | `cwe-251-300.json` |
| `CWE-281` | Improper Preservation of Permissions | C | Security | Medium | `cwe-251-300.json` |
| `CWE-283` | Unverified Ownership | C | General | Medium | `cwe-251-300.json` |
| `CWE-294` | Authentication Bypass by Capture-replay | C | Security | Medium | `cwe-251-300.json` |
| `CWE-301` | Reflection Attack in an Authentication Protocol | C | Security | Medium | `cwe-301-350.json` |
| `CWE-303` | Incorrect Implementation of Authentication Algorithm | C | Security | Medium | `cwe-301-350.json` |
| `CWE-305` | Authentication Bypass by Primary Weakness | C | Security | Medium | `cwe-301-350.json` |
| `CWE-308` | Use of Single-factor Authentication | C | Security | Medium | `cwe-301-350.json` |
| `CWE-309` | Use of Password System for Primary Authentication | C | Security | Medium | `cwe-301-350.json` |
| `CWE-322` | Key Exchange without Entity Authentication | C | Security | Medium | `cwe-301-350.json` |
| `CWE-324` | Use of a Key Past its Expiration Date | C | General | Medium | `cwe-301-350.json` |
| `CWE-325` | Missing Cryptographic Step | C | Cryptography | Medium | `cwe-301-350.json` |
| `CWE-341` | Predictable from Observable State | C | General | Medium | `cwe-301-350.json` |
| `CWE-342` | Predictable Exact Value from Previous Values | C | General | Medium | `cwe-301-350.json` |
| `CWE-343` | Predictable Value Range from Previous Values | C | General | Medium | `cwe-301-350.json` |
| `CWE-344` | Use of Invariant Value in Dynamically Changing Context | C | General | Medium | `cwe-301-350.json` |
| `CWE-348` | Use of Less Trusted Source | C | General | Medium | `cwe-301-350.json` |
| `CWE-349` | Acceptance of Extraneous Untrusted Data With Trusted Data | C | General | Medium | `cwe-301-350.json` |
| `CWE-351` | Insufficient Type Distinction | C | General | Medium | `cwe-351-400.json` |
| `CWE-353` | Missing Support for Integrity Check | C | General | Medium | `cwe-351-400.json` |
| `CWE-354` | Improper Validation of Integrity Check Value | C | General | Medium | `cwe-351-400.json` |
| `CWE-356` | Product UI does not Warn User of Unsafe Actions | C | General | Medium | `cwe-351-400.json` |
| `CWE-357` | Insufficient UI Warning of Dangerous Operations | C | General | Medium | `cwe-351-400.json` |
| `CWE-358` | Improperly Implemented Security Check for Standard | C | General | Medium | `cwe-351-400.json` |
| `CWE-360` | Trust of System Event Data | C | General | Medium | `cwe-351-400.json` |
| `CWE-368` | Context Switching Race Condition | C | Resource Management | Medium | `cwe-351-400.json` |
| `CWE-369` | Divide By Zero | C | Numeric | Medium | `cwe-351-400.json` |
| `CWE-372` | Incomplete Internal State Distinction | C | General | Medium | `cwe-351-400.json` |
| `CWE-385` | Covert Timing Channel | C | General | Medium | `cwe-351-400.json` |
| `CWE-386` | Symbolic Name not Mapping to Correct Object | C | General | Medium | `cwe-351-400.json` |
| `CWE-391` | Unchecked Error Condition | C | Error Handling | Medium | `cwe-351-400.json` |
| `CWE-392` | Missing Report of Error Condition | C | Error Handling | Medium | `cwe-351-400.json` |
| `CWE-393` | Return of Wrong Status Code | C | General | Medium | `cwe-351-400.json` |
| `CWE-394` | Unexpected Status Code or Return Value | C | General | Medium | `cwe-351-400.json` |
| `CWE-408` | Incorrect Behavior Order: Early Amplification | C | General | Medium | `cwe-401-450.json` |
| `CWE-410` | Insufficient Resource Pool | C | Resource Management | Medium | `cwe-401-450.json` |
| `CWE-412` | Unrestricted Externally Accessible Lock | C | Resource Management | Medium | `cwe-401-450.json` |
| `CWE-413` | Improper Resource Locking | C | Resource Management | Medium | `cwe-401-450.json` |
| `CWE-414` | Missing Lock Check | C | Resource Management | Medium | `cwe-401-450.json` |
| `CWE-419` | Unprotected Primary Channel | C | General | Medium | `cwe-401-450.json` |
| `CWE-420` | Unprotected Alternate Channel | C | General | Medium | `cwe-401-450.json` |
| `CWE-421` | Race Condition During Access to Alternate Channel | C | Program Invocation | Medium | `cwe-401-450.json` |
| `CWE-425` | Direct Request ('Forced Browsing') | C | General | Medium | `cwe-401-450.json` |
| `CWE-428` | Unquoted Search Path or Element | C | Program Invocation | Medium | `cwe-401-450.json` |
| `CWE-430` | Deployment of Wrong Handler | C | General | Medium | `cwe-401-450.json` |
| `CWE-431` | Missing Handler | C | General | Medium | `cwe-401-450.json` |
| `CWE-437` | Incomplete Model of Endpoint Features | C | General | Medium | `cwe-401-450.json` |
| `CWE-439` | Behavioral Change in New Version or Environment | C | General | Medium | `cwe-401-450.json` |
| `CWE-440` | Expected Behavior Violation | C | General | Medium | `cwe-401-450.json` |
| `CWE-444` | Inconsistent Interpretation of HTTP Requests ('HTTP Request/Response Smuggling') | C | General | Medium | `cwe-401-450.json` |
| `CWE-447` | Unimplemented or Unsupported Feature in UI | C | General | Medium | `cwe-401-450.json` |
| `CWE-448` | Obsolete Feature in UI | C | General | Medium | `cwe-401-450.json` |
| `CWE-449` | The UI Performs the Wrong Action | C | General | Medium | `cwe-401-450.json` |
| `CWE-455` | Non-exit on Failed Initialization | C | General | Medium | `cwe-451-500.json` |
| `CWE-474` | Use of Function with Inconsistent Implementations | C | General | Medium | `cwe-451-500.json` |
| `CWE-475` | Undefined Behavior for Input to API | C | General | Medium | `cwe-451-500.json` |
| `CWE-480` | Use of Incorrect Operator | C | General | Medium | `cwe-451-500.json` |
| `CWE-484` | Omitted Break Statement in Switch | C | General | Medium | `cwe-451-500.json` |
| `CWE-501` | Trust Boundary Violation | C | General | Medium | `cwe-501-550.json` |
| `CWE-515` | Covert Storage Channel | C | General | Medium | `cwe-501-550.json` |
| `CWE-544` | Missing Standardized Error Handling Mechanism | C | Error Handling | Medium | `cwe-501-550.json` |
| `CWE-549` | Missing Password Field Masking | C | Security | Medium | `cwe-501-550.json` |
| `CWE-551` | Incorrect Behavior Order: Authorization Before Parsing and Canonicalization | C | Security | Medium | `cwe-551-600.json` |
| `CWE-561` | Dead Code | C | General | Medium | `cwe-551-600.json` |
| `CWE-563` | Assignment to Variable without Use | C | General | Medium | `cwe-551-600.json` |
| `CWE-570` | Expression is Always False | C | General | Medium | `cwe-551-600.json` |
| `CWE-571` | Expression is Always True | C | General | Medium | `cwe-551-600.json` |
| `CWE-603` | Use of Client-Side Authentication | C | Security | Medium | `cwe-601-650.json` |
| `CWE-606` | Unchecked Input for Loop Condition | C | General | Medium | `cwe-601-650.json` |
| `CWE-618` | Exposed Unsafe ActiveX Method | C | General | Medium | `cwe-601-650.json` |
| `CWE-620` | Unverified Password Change | C | Security | Medium | `cwe-601-650.json` |
| `CWE-628` | Function Call with Incorrectly Specified Arguments | C | General | Medium | `cwe-601-650.json` |
| `CWE-639` | Authorization Bypass Through User-Controlled Key | C | Security | Medium | `cwe-601-650.json` |
| `CWE-640` | Weak Password Recovery Mechanism for Forgotten Password | C | Security | Medium | `cwe-601-650.json` |
| `CWE-645` | Overly Restrictive Account Lockout Mechanism | C | Resource Management | Medium | `cwe-601-650.json` |
| `CWE-649` | Reliance on Obfuscation or Encryption of Security-Relevant Inputs without Integrity Checking | C | General | Medium | `cwe-601-650.json` |
| `CWE-653` | Improper Isolation or Compartmentalization | C | General | Medium | `cwe-651-700.json` |
| `CWE-654` | Reliance on a Single Factor in a Security Decision | C | General | Medium | `cwe-651-700.json` |
| `CWE-656` | Reliance on Security Through Obscurity | C | General | Medium | `cwe-651-700.json` |
| `CWE-663` | Use of a Non-reentrant Function in a Concurrent Context | C | General | Medium | `cwe-651-700.json` |
| `CWE-681` | Incorrect Conversion between Numeric Types | C | Numeric | Medium | `cwe-651-700.json` |
| `CWE-694` | Use of Multiple Resources with Duplicate Identifier | C | Resource Management | Medium | `cwe-651-700.json` |
| `CWE-764` | Multiple Locks of a Critical Resource | C | Resource Management | Medium | `cwe-751-800.json` |
| `CWE-765` | Multiple Unlocks of a Critical Resource | C | Resource Management | Medium | `cwe-751-800.json` |
| `CWE-771` | Missing Reference to Active Allocated Resource | C | Resource Management | Medium | `cwe-751-800.json` |
| `CWE-778` | Insufficient Logging | C | Information Disclosure | Medium | `cwe-751-800.json` |
| `CWE-783` | Operator Precedence Logic Error | C | Information Disclosure | Medium | `cwe-751-800.json` |
| `CWE-791` | Incomplete Filtering of Special Elements | C | General | Medium | `cwe-751-800.json` |
| `CWE-804` | Guessable CAPTCHA | C | General | Medium | `cwe-801-850.json` |
| `CWE-820` | Missing Synchronization | C | General | Medium | `cwe-801-850.json` |
| `CWE-821` | Incorrect Synchronization | C | General | Medium | `cwe-801-850.json` |
| `CWE-826` | Premature Release of Resource During Expected Lifetime | C | Resource Management | Medium | `cwe-801-850.json` |
| `CWE-832` | Unlock of a Resource that is not Locked | C | Resource Management | Medium | `cwe-801-850.json` |
| `CWE-833` | Deadlock | C | Resource Management | Medium | `cwe-801-850.json` |
| `CWE-835` | Loop with Unreachable Exit Condition ('Infinite Loop') | C | General | Medium | `cwe-801-850.json` |
| `CWE-837` | Improper Enforcement of a Single, Unique Action | C | General | Medium | `cwe-801-850.json` |
| `CWE-841` | Improper Enforcement of Behavioral Workflow | C | General | Medium | `cwe-801-850.json` |
| `CWE-842` | Placement of User into Incorrect Group | C | General | Medium | `cwe-801-850.json` |
| `CWE-1024` | Comparison of Incompatible Types | C | Numeric | Medium | `cwe-1001-1050.json` |
| `CWE-1025` | Comparison Using Wrong Factors | C | Numeric | Medium | `cwe-1001-1050.json` |
| `CWE-1037` | Processor Optimization Removal or Modification of Security-critical Code | C | General | Medium | `cwe-1001-1050.json` |
| `CWE-1041` | Use of Redundant Code | C | General | Medium | `cwe-1001-1050.json` |
| `CWE-1043` | Data Element Aggregating an Excessively Large Number of Non-Primitive Elements | C | General | Medium | `cwe-1001-1050.json` |
| `CWE-1044` | Architecture with Number of Horizontal Layers Outside of Expected Range | C | General | Medium | `cwe-1001-1050.json` |
| `CWE-1045` | Parent Class with a Virtual Destructor and a Child Class without a Virtual Destructor | C | General | Medium | `cwe-1001-1050.json` |
| `CWE-1047` | Modules with Circular Dependencies | C | General | Medium | `cwe-1001-1050.json` |
| `CWE-1048` | Invokable Control Element with Large Number of Outward Calls | C | General | Medium | `cwe-1001-1050.json` |
| `CWE-1049` | Excessive Data Query Operations in a Large Data Table | C | General | Medium | `cwe-1001-1050.json` |
| `CWE-1053` | Missing Documentation for Design | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1054` | Invocation of a Control Element at an Unnecessarily Deep Horizontal Layer | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1055` | Multiple Inheritance from Concrete Classes | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1056` | Invokable Control Element with Variadic Parameters | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1057` | Data Access Operations Outside of Expected Data Manager Component | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1058` | Invokable Control Element in Multi-Thread Context with non-Final Static Storable or Member Element | C | Resource Management | Medium | `cwe-1051-1100.json` |
| `CWE-1062` | Parent Class with References to Child Class | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1063` | Creation of Class Instance within a Static Code Block | C | Resource Management | Medium | `cwe-1051-1100.json` |
| `CWE-1064` | Invokable Control Element with Signature Containing an Excessive Number of Parameters | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1065` | Runtime Resource Management Control Element in a Component Built to Run on Application Servers | C | Resource Management | Medium | `cwe-1051-1100.json` |
| `CWE-1066` | Missing Serialization Control Element | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1068` | Inconsistency Between Implementation and Documented Design | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1070` | Serializable Data Element Containing non-Serializable Item Elements | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1074` | Class with Excessively Deep Inheritance | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1075` | Unconditional Control Flow Transfer outside of Switch Block | C | Resource Management | Medium | `cwe-1051-1100.json` |
| `CWE-1079` | Parent Class without Virtual Destructor Method | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1080` | Source Code File with Excessive Number of Lines of Code | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1082` | Class Instance Self Destruction Control Element | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1083` | Data Access from Outside Expected Data Manager Component | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1085` | Invokable Control Element with Excessive Volume of Commented-out Code | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1086` | Class with Excessive Number of Child Classes | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1087` | Class with Virtual Method without a Virtual Destructor | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1089` | Large Data Table with Excessive Number of Indices | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1090` | Method Containing Access of a Member Element from Another Class | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1092` | Use of Same Invokable Control Element in Multiple Architectural Layers | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1094` | Excessive Index Range Scan for a Data Resource | C | Resource Management | Medium | `cwe-1051-1100.json` |
| `CWE-1095` | Loop Condition Value Update within the Loop | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1097` | Persistent Storable Data Element without Associated Comparison Control Element | C | Numeric | Medium | `cwe-1051-1100.json` |
| `CWE-1098` | Data Element containing Pointer Item without Proper Copy Control Element | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1099` | Inconsistent Naming Conventions for Identifiers | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1100` | Insufficient Isolation of System-Dependent Functions | C | General | Medium | `cwe-1051-1100.json` |
| `CWE-1101` | Reliance on Runtime Component in Generated Code | C | General | Medium | `cwe-1101-1150.json` |
| `CWE-1102` | Reliance on Machine-Dependent Data Representation | C | General | Medium | `cwe-1101-1150.json` |
| `CWE-1103` | Use of Platform-Dependent Third Party Components | C | General | Medium | `cwe-1101-1150.json` |
| `CWE-1105` | Insufficient Encapsulation of Machine-Dependent Functionality | C | General | Medium | `cwe-1101-1150.json` |
| `CWE-1107` | Insufficient Isolation of Symbolic Constant Definitions | C | General | Medium | `cwe-1101-1150.json` |
| `CWE-1109` | Use of Same Variable for Multiple Purposes | C | General | Medium | `cwe-1101-1150.json` |
| `CWE-1110` | Incomplete Design Documentation | C | General | Medium | `cwe-1101-1150.json` |
| `CWE-1111` | Incomplete I/O Documentation | C | General | Medium | `cwe-1101-1150.json` |
| `CWE-1112` | Incomplete Documentation of Program Execution | C | General | Medium | `cwe-1101-1150.json` |
| `CWE-1113` | Inappropriate Comment Style | C | General | Medium | `cwe-1101-1150.json` |
| `CWE-1114` | Inappropriate Whitespace Style | C | General | Medium | `cwe-1101-1150.json` |
| `CWE-1115` | Source Code Element without Standard Prologue | C | Information Disclosure | Medium | `cwe-1101-1150.json` |
| `CWE-1116` | Inaccurate Source Code Comments | C | General | Medium | `cwe-1101-1150.json` |
| `CWE-1117` | Callable with Insufficient Behavioral Summary | C | General | Medium | `cwe-1101-1150.json` |
| `CWE-1118` | Insufficient Documentation of Error Handling Techniques | C | Error Handling | Medium | `cwe-1101-1150.json` |
| `CWE-1119` | Excessive Use of Unconditional Branching | C | General | Medium | `cwe-1101-1150.json` |
| `CWE-1122` | Excessive Halstead Complexity | C | General | Medium | `cwe-1101-1150.json` |
| `CWE-1126` | Declaration of Variable with Unnecessarily Wide Scope | C | General | Medium | `cwe-1101-1150.json` |
| `CWE-1127` | Compilation with Insufficient Warnings or Errors | C | Error Handling | Medium | `cwe-1101-1150.json` |

## Executive Summary

These catalogue entries remain **registered in JSON only**. They lack a high-signal pure-Go source-pattern detector for v0 (design-level, process, CISQ maintainability, needs whole-program analysis, or no stable Python sink token).

**Do not** open implement PRs against this list without promoting an ID to a themed batch (update ownership + inventory).

## Promotion path

1. Identify FN-safe needles + hit/miss fixtures.
2. Move ID from this file into an implement batch (or new batch-17+).
3. Update `_inventory.json` deferred → missing and tier C → A/B.

## Deferred catalogue (ledger)

- [~] `CWE-115` — Misinterpretation of Input (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-134` — Use of Externally-Controlled Format String (Logging) — deferred: no v0 pure-Go sink plan
- [~] `CWE-166` — Improper Handling of Missing Special Element (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-167` — Improper Handling of Additional Special Element (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-168` — Improper Handling of Inconsistent Special Elements (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-183` — Permissive List of Allowed Inputs (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-190` — Integer Overflow or Wraparound (Number Processing) — deferred: no v0 pure-Go sink plan
- [~] `CWE-193` — Off-by-one Error (Error Handling) — deferred: no v0 pure-Go sink plan
- [~] `CWE-205` — Observable Behavioral Discrepancy (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-222` — Truncation of Security-relevant Information (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-223` — Omission of Security-relevant Information (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-224` — Obscured Security-relevant Information by Alternate Name (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-229` — Improper Handling of Values (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-233` — Improper Handling of Parameters (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-237` — Improper Handling of Structural Elements (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-241` — Improper Handling of Unexpected Data Type (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-253` — Incorrect Check of Function Return Value (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-262` — Not Using Password Aging (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-263` — Password Aging with Long Expiration (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-266` — Incorrect Privilege Assignment (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-267` — Privilege Defined With Unsafe Actions (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-268` — Privilege Chaining (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-270` — Privilege Context Switching Error (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-273` — Improper Check for Dropped Privileges (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-274` — Improper Handling of Insufficient Privileges (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-277` — Insecure Inherited Permissions (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-278` — Insecure Preserved Inherited Permissions (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-280` — Improper Handling of Insufficient Permissions or Privileges (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-281` — Improper Preservation of Permissions (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-283` — Unverified Ownership (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-294` — Authentication Bypass by Capture-replay (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-301` — Reflection Attack in an Authentication Protocol (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-303` — Incorrect Implementation of Authentication Algorithm (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-305` — Authentication Bypass by Primary Weakness (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-308` — Use of Single-factor Authentication (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-309` — Use of Password System for Primary Authentication (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-322` — Key Exchange without Entity Authentication (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-324` — Use of a Key Past its Expiration Date (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-325` — Missing Cryptographic Step (Cryptography) — deferred: no v0 pure-Go sink plan
- [~] `CWE-341` — Predictable from Observable State (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-342` — Predictable Exact Value from Previous Values (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-343` — Predictable Value Range from Previous Values (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-344` — Use of Invariant Value in Dynamically Changing Context (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-348` — Use of Less Trusted Source (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-349` — Acceptance of Extraneous Untrusted Data With Trusted Data (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-351` — Insufficient Type Distinction (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-353` — Missing Support for Integrity Check (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-354` — Improper Validation of Integrity Check Value (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-356` — Product UI does not Warn User of Unsafe Actions (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-357` — Insufficient UI Warning of Dangerous Operations (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-358` — Improperly Implemented Security Check for Standard (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-360` — Trust of System Event Data (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-368` — Context Switching Race Condition (Resource Management) — deferred: no v0 pure-Go sink plan
- [~] `CWE-369` — Divide By Zero (Numeric) — deferred: no v0 pure-Go sink plan
- [~] `CWE-372` — Incomplete Internal State Distinction (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-385` — Covert Timing Channel (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-386` — Symbolic Name not Mapping to Correct Object (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-391` — Unchecked Error Condition (Error Handling) — deferred: no v0 pure-Go sink plan
- [~] `CWE-392` — Missing Report of Error Condition (Error Handling) — deferred: no v0 pure-Go sink plan
- [~] `CWE-393` — Return of Wrong Status Code (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-394` — Unexpected Status Code or Return Value (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-408` — Incorrect Behavior Order: Early Amplification (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-410` — Insufficient Resource Pool (Resource Management) — deferred: no v0 pure-Go sink plan
- [~] `CWE-412` — Unrestricted Externally Accessible Lock (Resource Management) — deferred: no v0 pure-Go sink plan
- [~] `CWE-413` — Improper Resource Locking (Resource Management) — deferred: no v0 pure-Go sink plan
- [~] `CWE-414` — Missing Lock Check (Resource Management) — deferred: no v0 pure-Go sink plan
- [~] `CWE-419` — Unprotected Primary Channel (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-420` — Unprotected Alternate Channel (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-421` — Race Condition During Access to Alternate Channel (Program Invocation) — deferred: no v0 pure-Go sink plan
- [~] `CWE-425` — Direct Request ('Forced Browsing') (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-428` — Unquoted Search Path or Element (Program Invocation) — deferred: no v0 pure-Go sink plan
- [~] `CWE-430` — Deployment of Wrong Handler (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-431` — Missing Handler (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-437` — Incomplete Model of Endpoint Features (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-439` — Behavioral Change in New Version or Environment (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-440` — Expected Behavior Violation (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-444` — Inconsistent Interpretation of HTTP Requests ('HTTP Request/Response Smuggling') (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-447` — Unimplemented or Unsupported Feature in UI (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-448` — Obsolete Feature in UI (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-449` — The UI Performs the Wrong Action (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-455` — Non-exit on Failed Initialization (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-474` — Use of Function with Inconsistent Implementations (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-475` — Undefined Behavior for Input to API (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-480` — Use of Incorrect Operator (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-484` — Omitted Break Statement in Switch (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-501` — Trust Boundary Violation (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-515` — Covert Storage Channel (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-544` — Missing Standardized Error Handling Mechanism (Error Handling) — deferred: no v0 pure-Go sink plan
- [~] `CWE-549` — Missing Password Field Masking (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-551` — Incorrect Behavior Order: Authorization Before Parsing and Canonicalization (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-561` — Dead Code (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-563` — Assignment to Variable without Use (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-570` — Expression is Always False (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-571` — Expression is Always True (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-603` — Use of Client-Side Authentication (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-606` — Unchecked Input for Loop Condition (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-618` — Exposed Unsafe ActiveX Method (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-620` — Unverified Password Change (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-628` — Function Call with Incorrectly Specified Arguments (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-639` — Authorization Bypass Through User-Controlled Key (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-640` — Weak Password Recovery Mechanism for Forgotten Password (Security) — deferred: no v0 pure-Go sink plan
- [~] `CWE-645` — Overly Restrictive Account Lockout Mechanism (Resource Management) — deferred: no v0 pure-Go sink plan
- [~] `CWE-649` — Reliance on Obfuscation or Encryption of Security-Relevant Inputs without Integrity Checking (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-653` — Improper Isolation or Compartmentalization (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-654` — Reliance on a Single Factor in a Security Decision (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-656` — Reliance on Security Through Obscurity (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-663` — Use of a Non-reentrant Function in a Concurrent Context (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-681` — Incorrect Conversion between Numeric Types (Numeric) — deferred: no v0 pure-Go sink plan
- [~] `CWE-694` — Use of Multiple Resources with Duplicate Identifier (Resource Management) — deferred: no v0 pure-Go sink plan
- [~] `CWE-764` — Multiple Locks of a Critical Resource (Resource Management) — deferred: no v0 pure-Go sink plan
- [~] `CWE-765` — Multiple Unlocks of a Critical Resource (Resource Management) — deferred: no v0 pure-Go sink plan
- [~] `CWE-771` — Missing Reference to Active Allocated Resource (Resource Management) — deferred: no v0 pure-Go sink plan
- [~] `CWE-778` — Insufficient Logging (Information Disclosure) — deferred: no v0 pure-Go sink plan
- [~] `CWE-783` — Operator Precedence Logic Error (Information Disclosure) — deferred: no v0 pure-Go sink plan
- [~] `CWE-791` — Incomplete Filtering of Special Elements (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-804` — Guessable CAPTCHA (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-820` — Missing Synchronization (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-821` — Incorrect Synchronization (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-826` — Premature Release of Resource During Expected Lifetime (Resource Management) — deferred: no v0 pure-Go sink plan
- [~] `CWE-832` — Unlock of a Resource that is not Locked (Resource Management) — deferred: no v0 pure-Go sink plan
- [~] `CWE-833` — Deadlock (Resource Management) — deferred: no v0 pure-Go sink plan
- [~] `CWE-835` — Loop with Unreachable Exit Condition ('Infinite Loop') (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-837` — Improper Enforcement of a Single, Unique Action (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-841` — Improper Enforcement of Behavioral Workflow (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-842` — Placement of User into Incorrect Group (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1024` — Comparison of Incompatible Types (Numeric) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1025` — Comparison Using Wrong Factors (Numeric) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1037` — Processor Optimization Removal or Modification of Security-critical Code (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1041` — Use of Redundant Code (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1043` — Data Element Aggregating an Excessively Large Number of Non-Primitive Elements (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1044` — Architecture with Number of Horizontal Layers Outside of Expected Range (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1045` — Parent Class with a Virtual Destructor and a Child Class without a Virtual Destructor (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1047` — Modules with Circular Dependencies (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1048` — Invokable Control Element with Large Number of Outward Calls (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1049` — Excessive Data Query Operations in a Large Data Table (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1053` — Missing Documentation for Design (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1054` — Invocation of a Control Element at an Unnecessarily Deep Horizontal Layer (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1055` — Multiple Inheritance from Concrete Classes (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1056` — Invokable Control Element with Variadic Parameters (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1057` — Data Access Operations Outside of Expected Data Manager Component (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1058` — Invokable Control Element in Multi-Thread Context with non-Final Static Storable or Member Element (Resource Management) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1062` — Parent Class with References to Child Class (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1063` — Creation of Class Instance within a Static Code Block (Resource Management) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1064` — Invokable Control Element with Signature Containing an Excessive Number of Parameters (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1065` — Runtime Resource Management Control Element in a Component Built to Run on Application Servers (Resource Management) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1066` — Missing Serialization Control Element (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1068` — Inconsistency Between Implementation and Documented Design (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1070` — Serializable Data Element Containing non-Serializable Item Elements (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1074` — Class with Excessively Deep Inheritance (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1075` — Unconditional Control Flow Transfer outside of Switch Block (Resource Management) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1079` — Parent Class without Virtual Destructor Method (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1080` — Source Code File with Excessive Number of Lines of Code (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1082` — Class Instance Self Destruction Control Element (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1083` — Data Access from Outside Expected Data Manager Component (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1085` — Invokable Control Element with Excessive Volume of Commented-out Code (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1086` — Class with Excessive Number of Child Classes (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1087` — Class with Virtual Method without a Virtual Destructor (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1089` — Large Data Table with Excessive Number of Indices (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1090` — Method Containing Access of a Member Element from Another Class (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1092` — Use of Same Invokable Control Element in Multiple Architectural Layers (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1094` — Excessive Index Range Scan for a Data Resource (Resource Management) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1095` — Loop Condition Value Update within the Loop (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1097` — Persistent Storable Data Element without Associated Comparison Control Element (Numeric) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1098` — Data Element containing Pointer Item without Proper Copy Control Element (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1099` — Inconsistent Naming Conventions for Identifiers (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1100` — Insufficient Isolation of System-Dependent Functions (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1101` — Reliance on Runtime Component in Generated Code (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1102` — Reliance on Machine-Dependent Data Representation (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1103` — Use of Platform-Dependent Third Party Components (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1105` — Insufficient Encapsulation of Machine-Dependent Functionality (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1107` — Insufficient Isolation of Symbolic Constant Definitions (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1109` — Use of Same Variable for Multiple Purposes (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1110` — Incomplete Design Documentation (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1111` — Incomplete I/O Documentation (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1112` — Incomplete Documentation of Program Execution (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1113` — Inappropriate Comment Style (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1114` — Inappropriate Whitespace Style (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1115` — Source Code Element without Standard Prologue (Information Disclosure) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1116` — Inaccurate Source Code Comments (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1117` — Callable with Insufficient Behavioral Summary (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1118` — Insufficient Documentation of Error Handling Techniques (Error Handling) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1119` — Excessive Use of Unconditional Branching (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1122` — Excessive Halstead Complexity (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1126` — Declaration of Variable with Unnecessarily Wide Scope (General) — deferred: no v0 pure-Go sink plan
- [~] `CWE-1127` — Compilation with Insufficient Warnings or Errors (Error Handling) — deferred: no v0 pure-Go sink plan

**Count:** 185
