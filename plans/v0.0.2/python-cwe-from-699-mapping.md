# Python CWE mapping from CWE view 699.csv

> **Source:** `699.csv` (CWE-699 Software Development view export)
> **Rows in CSV:** 399
> **Included for Python catalogue:** 344
> **Excluded:** 55
> **Max rules per JSON file:** 50
> **Output:** `ruleset/python/chunks/cwe-*.json`

## Filter rules

| Decision | Rule |
|----------|------|
| **Include · python-specific** | `Applicable Platforms` lists `LANGUAGE NAME:Python` |
| **Include · generic** | `LANGUAGE CLASS:Not Language-Specific` (or no exclusive non-Python language list) |
| **Exclude · lang_only** | Platforms name other languages only (C/C++, Java, PHP, …) without Not Language-Specific / Python |
| **Exclude · memory_corruption_c_style** | C-style memory safety weakness names (buffer overflow, UAF, double free, …) |
| **Exclude · memory_unsafe_class** | `LANGUAGE CLASS:Memory-Unsafe` without Not Language-Specific / Python |

## Summary counts

- **Generic (portable):** 339
- **Python-specific platform:** 5
- **Excluded:** 55

## Python-specific platform CWEs

| CWE | Name |
|-----|------|
| CWE-396 | Declaration of Catch for Generic Exception |
| CWE-397 | Declaration of Throws for Generic Exception |
| CWE-478 | Missing Default Case in Multiple Condition Expression |
| CWE-502 | Deserialization of Untrusted Data |
| CWE-915 | Improperly Controlled Modification of Dynamically-Determined Object Attributes |

## Output chunk files (≤50 rules each)

| File | Count | ID min | ID max |
|------|------:|-------:|-------:|
| `cwe-001-050.json` | 3 | 15 | 41 |
| `cwe-051-100.json` | 12 | 59 | 94 |
| `cwe-101-150.json` | 5 | 112 | 140 |
| `cwe-151-200.json` | 11 | 166 | 193 |
| `cwe-201-250.json` | 17 | 201 | 250 |
| `cwe-251-300.json` | 26 | 252 | 295 |
| `cwe-301-350.json` | 26 | 301 | 349 |
| `cwe-351-400.json` | 23 | 351 | 397 |
| `cwe-401-450.json` | 24 | 403 | 449 |
| `cwe-451-500.json` | 14 | 454 | 497 |
| `cwe-501-550.json` | 10 | 501 | 549 |
| `cwe-551-600.json` | 8 | 551 | 584 |
| `cwe-601-650.json` | 16 | 601 | 649 |
| `cwe-651-700.json` | 8 | 653 | 698 |
| `cwe-701-750.json` | 2 | 708 | 749 |
| `cwe-751-800.json` | 12 | 756 | 798 |
| `cwe-801-850.json` | 14 | 804 | 842 |
| `cwe-901-950.json` | 14 | 908 | 941 |
| `cwe-1001-1050.json` | 14 | 1007 | 1050 |
| `cwe-1051-1100.json` | 38 | 1051 | 1100 |
| `cwe-1101-1150.json` | 26 | 1101 | 1127 |
| `cwe-1151-1200.json` | 2 | 1173 | 1188 |
| `cwe-1201-1250.json` | 6 | 1204 | 1241 |
| `cwe-1251-1300.json` | 7 | 1265 | 1289 |
| `cwe-1301-1350.json` | 4 | 1322 | 1341 |
| `cwe-1351-1400.json` | 2 | 1389 | 1392 |

## Excluded by reason

### lang_only:C,C++ (30)

- CWE-120: Buffer Copy without Checking Size of Input ('Classic Buffer Overflow')
- CWE-124: Buffer Underwrite ('Buffer Underflow')
- CWE-125: Out-of-bounds Read
- CWE-128: Wrap-around Error
- CWE-131: Incorrect Calculation of Buffer Size
- CWE-135: Incorrect Calculation of Multi-Byte String Length
- CWE-170: Improper Null Termination
- CWE-242: Use of Inherently Dangerous Function
- CWE-243: Creation of chroot Jail Without Changing Working Directory
- CWE-364: Signal Handler Race Condition
- CWE-463: Deletion of Data Structure Sentinel
- CWE-464: Addition of Data Structure Sentinel
- CWE-466: Return of Pointer Value Outside of Expected Range
- CWE-468: Incorrect Pointer Scaling
- CWE-469: Use of Pointer Subtraction to Determine Size
- CWE-483: Incorrect Block Delimitation
- CWE-562: Return of Stack Variable Address
- CWE-676: Use of Potentially Dangerous Function
- CWE-733: Compiler Optimization Removal or Modification of Security-critical Code
- CWE-763: Release of Invalid Pointer or Reference
- … +10 more

### lang_only:Java (7)

- CWE-395: Use of NullPointerException Catch to Detect NULL Pointer Dereference
- CWE-487: Reliance on Package-level Scope
- CWE-580: clone() Method Without super.clone()
- CWE-581: Object Model Violation: Just One of Equals and Hashcode Defined
- CWE-586: Explicit Call to Finalize()
- CWE-609: Double-Checked Locking
- CWE-917: Improper Neutralization of Special Elements used in an Expression Language Statement ('Expression Language Injection')

### lang_only:C,C#,C++,Java (4)

- CWE-191: Integer Underflow (Wrap or Wraparound)
- CWE-366: Race Condition within a Thread
- CWE-374: Passing Mutable Objects to an Untrusted Method
- CWE-375: Returning a Mutable Object to an Untrusted Caller

### lang_only:C#,C++,Java (3)

- CWE-248: Uncaught Exception
- CWE-766: Critical Data Element Declared Public
- CWE-767: Access to Critical Private Variable via Public Method

### lang_only:SQL (2)

- CWE-619: Dangling Database Cursor ('Cursor Injection')
- CWE-1073: Non-SQL Invokable Control Element with Excessive Number of Data Resource Accesses

### lang_only:PHP,Perl (2)

- CWE-624: Executable Regular Expression Error
- CWE-625: Permissive Regular Expression

### memory_corruption_c_style (1)

- CWE-130: Improper Handling of Length Parameter Inconsistency

### lang_only:Java,PHP (1)

- CWE-470: Use of Externally-Controlled Input to Select Classes or Code ('Unsafe Reflection')

### lang_only:C,C#,C++,Go,Java (1)

- CWE-476: NULL Pointer Dereference

### lang_only:C,C#,C++ (1)

- CWE-587: Assignment of a Fixed Address to a Pointer

### lang_only:C#,Java (1)

- CWE-1235: Incorrect Use of Autoboxing and Unboxing for Performance Critical Operations

### lang_only:Other (1)

- CWE-1327: Binding to an Unrestricted IP Address

### lang_only:C,C#,C++,Java,JavaScript (1)

- CWE-1335: Incorrect Bitwise Shift of Integer

## Full included ID list

CWE-15, CWE-22, CWE-41, CWE-59, CWE-66, CWE-73, CWE-76, CWE-78, CWE-79, CWE-88, CWE-89, CWE-90, CWE-91, CWE-93, CWE-94, CWE-112, CWE-115, CWE-117, CWE-134, CWE-140, CWE-166, CWE-167, CWE-168, CWE-178, CWE-179, CWE-182, CWE-183, CWE-184, CWE-186, CWE-190, CWE-193, CWE-201, CWE-204, CWE-205, CWE-208, CWE-209, CWE-212, CWE-213, CWE-214, CWE-215, CWE-222, CWE-223, CWE-224, CWE-229, CWE-233, CWE-237, CWE-241, CWE-250, CWE-252, CWE-253, CWE-256, CWE-257, CWE-260, CWE-261, CWE-262, CWE-263, CWE-266, CWE-267, CWE-268, CWE-270, CWE-272, CWE-273, CWE-274, CWE-276, CWE-277, CWE-278, CWE-279, CWE-280, CWE-281, CWE-283, CWE-289, CWE-290, CWE-294, CWE-295, CWE-301, CWE-303, CWE-305, CWE-306, CWE-307, CWE-308, CWE-309, CWE-312, CWE-319, CWE-322, CWE-323, CWE-324, CWE-325, CWE-328, CWE-331, CWE-334, CWE-335, CWE-338, CWE-341, CWE-342, CWE-343, CWE-344, CWE-346, CWE-347, CWE-348, CWE-349, CWE-351, CWE-353, CWE-354, CWE-356, CWE-357, CWE-358, CWE-359, CWE-360, CWE-367, CWE-368, CWE-369, CWE-372, CWE-378, CWE-379, CWE-385, CWE-386, CWE-390, CWE-391, CWE-392, CWE-393, CWE-394, CWE-396, CWE-397, CWE-403, CWE-408, CWE-409, CWE-410, CWE-412, CWE-413, CWE-414, CWE-419, CWE-420, CWE-421, CWE-425, CWE-426, CWE-427, CWE-428, CWE-430, CWE-431, CWE-434, CWE-437, CWE-439, CWE-440, CWE-444, CWE-447, CWE-448, CWE-449, CWE-454, CWE-455, CWE-459, CWE-472, CWE-474, CWE-475, CWE-477, CWE-478, CWE-480, CWE-484, CWE-488, CWE-489, CWE-494, CWE-497, CWE-501, CWE-502, CWE-515, CWE-521, CWE-523, CWE-524, CWE-538, CWE-544, CWE-547, CWE-549, CWE-551, CWE-552, CWE-561, CWE-563, CWE-565, CWE-570, CWE-571, CWE-584, CWE-601, CWE-603, CWE-605, CWE-606, CWE-611, CWE-613, CWE-617, CWE-618, CWE-620, CWE-628, CWE-639, CWE-640, CWE-641, CWE-645, CWE-648, CWE-649, CWE-653, CWE-654, CWE-656, CWE-663, CWE-681, CWE-694, CWE-695, CWE-698, CWE-708, CWE-749, CWE-756, CWE-764, CWE-765, CWE-770, CWE-771, CWE-772, CWE-776, CWE-778, CWE-779, CWE-783, CWE-791, CWE-798, CWE-804, CWE-807, CWE-820, CWE-821, CWE-826, CWE-829, CWE-832, CWE-833, CWE-835, CWE-836, CWE-837, CWE-838, CWE-841, CWE-842, CWE-908, CWE-909, CWE-910, CWE-911, CWE-914, CWE-915, CWE-916, CWE-918, CWE-920, CWE-921, CWE-924, CWE-939, CWE-940, CWE-941, CWE-1007, CWE-1021, CWE-1024, CWE-1025, CWE-1037, CWE-1041, CWE-1043, CWE-1044, CWE-1045, CWE-1046, CWE-1047, CWE-1048, CWE-1049, CWE-1050, CWE-1051, CWE-1052, CWE-1053, CWE-1054, CWE-1055, CWE-1056, CWE-1057, CWE-1058, CWE-1060, CWE-1062, CWE-1063, CWE-1064, CWE-1065, CWE-1066, CWE-1067, CWE-1068, CWE-1070, CWE-1071, CWE-1072, CWE-1074, CWE-1075, CWE-1079, CWE-1080, CWE-1082, CWE-1083, CWE-1084, CWE-1085, CWE-1086, CWE-1087, CWE-1089, CWE-1090, CWE-1092, CWE-1094, CWE-1095, CWE-1097, CWE-1098, CWE-1099, CWE-1100, CWE-1101, CWE-1102, CWE-1103, CWE-1104, CWE-1105, CWE-1106, CWE-1107, CWE-1108, CWE-1109, CWE-1110, CWE-1111, CWE-1112, CWE-1113, CWE-1114, CWE-1115, CWE-1116, CWE-1117, CWE-1118, CWE-1119, CWE-1121, CWE-1122, CWE-1123, CWE-1124, CWE-1125, CWE-1126, CWE-1127, CWE-1173, CWE-1188, CWE-1204, CWE-1220, CWE-1230, CWE-1236, CWE-1240, CWE-1241, CWE-1265, CWE-1284, CWE-1285, CWE-1286, CWE-1287, CWE-1288, CWE-1289, CWE-1322, CWE-1333, CWE-1339, CWE-1341, CWE-1389, CWE-1392

