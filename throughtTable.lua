function findCharacterTable(charID)
    for index, value in pairs(table["Default"]) do
        if type(value) == "table" then
            local charactersTable = table["Default"][index]["$AccountWide"]["characters"]
			for i, v in pairs(charactersTable) do
				if(i == charID) then
					print("Motif table for ", charID, " found")
					motifTable = charactersTable[i]["motifs"]
                else 
                    motifTable = "Character with ID: " .. charID .. " not found"
                end
            end
        end
	end
end