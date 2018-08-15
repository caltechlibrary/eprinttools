#
# issn_to_publisher_map.py generates a Go file for clsrules
# that creates a two map[string]string based on Journals 
# discovered in CaltechAUTHORS.  Journals.ds needs to 
# exists and be populated.
#
import dataset
import os

c = "Journals.ds"
keys = dataset.keys(c)
print("package clsrules")
print("")
print("var (")

# Generate a ISSN to Publisher Map
print("issnPublisher = map[string]string{")
for key in keys:
    try: 
        rec, err = dataset.read(c, key)
    except:
        rec = { "_Key": key, "publisher":"" }
    if err != "":
        print(f"// ERROR ({key}): {err}")
    print(f"    \"{rec['_Key']}\":\"{rec['publisher']}\",")
print("}")
print("")

# Generate a ISSN to Publication Map
print("issnPublication = map[string]string{")
for key in keys:
    try: 
        rec, err = dataset.read(c, key)
    except:
        rec = { "_Key": key, "publication":"" }
    if err != "":
        print(f"// ERROR ({key}): {err}")
    print(f"    \"{rec['_Key']}\":\"{rec['publication']}\",")
print("}")

print("")
print(")")
print("")
