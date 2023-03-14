package ngordnet.main;

import ngordnet.hugbrowsermagic.NgordnetQuery;
import ngordnet.hugbrowsermagic.NgordnetQueryHandler;
import ngordnet.ngrams.NGramMap;
import ngordnet.ngrams.TimeSeries;

import java.util.HashMap;
import java.util.List;

/* The constructor for HistoryTextHandler should be of the following form:
    public HistoryTextHandler(NGramMap map). */

public class HistoryTextHandler extends NgordnetQueryHandler {

    HashMap<String, TimeSeries> wordsMap;
    int x;
    TimeSeries<Double> wordsCountsData;
    NGramMap mapMap;

    public HistoryTextHandler (NGramMap map) {
        /* Need to access the data in NGramMap */
        wordsMap = map.returnNGramMapWord();
        mapMap = map;
        wordsCountsData = map.totalCountHistory();
    }

    @Override
    public String handle (NgordnetQuery q) {
        List<String> words = q.words();
        int startYear = q.startYear();
        int endYear = q.endYear();

        String response = " ";
        String curWord = " ";

        for (int i = 0; i < q.words().size(); i++) {
            curWord = q.words().get(i);
            response += curWord + ": ";
            response += mapMap.countHistory(curWord, startYear, endYear);
            response += "\n";
        }

        return response;
    }
}
