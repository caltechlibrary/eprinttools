#!/bin/bash
#
# Check for missing feed file types.

for FNAME in advisor-bachelors.json advisor-combined.json advisor-engd.json advisor-masters.json advisor-other.json advisor-phd.json advisor-senior_major.json advisor-senior_minor.json advisor.json article.json audiovisual.json bachelors.json book.json book_section.json caltechauthors-grid.json caltechdata-grid.json caltechthesis-grid.json collection.json combined.json combined_data.json conference_item.json data.json data_object_types.json data_pub_types.json data_types.json dataset.json directory_info.json engd.json group.json group_list.json image.json index.json interactiveresource.json masters.json model.json monograph.json object_types.json pagefind-entry.json patent.json people.json people_list.json phd.json pub_types.json senior_major.json senior_minor.json software.json teaching_resource.json text.json thesis.json video.json workflow.json; do
	R=$(find demo/htdocs -name "${FNAME}")
	if [ "$R" = "" ]; then
		echo "    - [ ] ${FNAME}"
	else
		echo "    - [x] ${FNAME}"
    fi
done
