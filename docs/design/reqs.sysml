package reqs{
    private import RequirementDerivation::*;
    requirement req1 {
        id = "REQ-001";
        text = "The system shall support user authentication.";
    }

    requirement req2 {
        id = "REQ-002";
        text = "The system shall provide data encryption for sensitive information.";
    }
    #derivation connection {
        end #original ::> req1;
        end #derive ::> req2;
    }
    //#derive Connector req1 to req2;
}