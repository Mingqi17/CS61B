package ngordnet.ngrams;

import edu.princeton.cs.algs4.In;

import java.sql.Time;
import java.util.*;

/** An object that provides utility methods for making queries on the
 *  Google NGrams dataset (or a subset thereof).
 *
 *  An NGramMap stores pertinent data from a "words file" and a "counts
 *  file". It is not a map n the strict sense, but it does provide additional
 *  functionality.
 *
 *  @author Josh Hug
 */
public class NGramMap {
    /** Constructs an NGramMap from WORDSFILENAME and COUNTSFILENAME. */
    HashMap<String, TimeSeries> totalMap;
    TimeSeries countsYear;
    public NGramMap(String wordsFilename, String countsFilename) {
        totalMap = new HashMap<String, TimeSeries>();
        In wordsFile = new In(wordsFilename);
        while (wordsFile.hasNextLine()) {
            String line = wordsFile.readLine();
            String[] splitLine = line.split("\t");
            String word = splitLine[0];
            Integer year = Integer.parseInt(splitLine[1]);
            Double numberOfTimes = Double.parseDouble(splitLine[2]);

            if (totalMap.containsKey(word)) {
                /* Does the totalMap.words key be updated actually? */
                TimeSeries<Number> ts = totalMap.get(word);
                ts.put(year, numberOfTimes.doubleValue());
            } else {
                TimeSeries<Number> ts = new TimeSeries<Number>();
                ts.put(year, numberOfTimes.doubleValue());
                totalMap.put(word, ts);
            }
        }

        countsYear = new TimeSeries<Number>();
        In countsFile = new In(countsFilename);
        while (countsFile.hasNextLine()) {
            String line = countsFile.readLine();
            String[] splitLine = line.split(",");
            Integer year = Integer.parseInt(splitLine[0]);
            Double numberOfWords = Double.parseDouble(splitLine[1]);
            countsYear.put(year, numberOfWords.doubleValue());
        }
    }

    public HashMap returnNGramMapWord () {
        return totalMap;
    }

    /** Provides the history of WORD. The returned TimeSeries should be a copy,
     *  not a link to this NGramMap's TimeSeries. In other words, changes made
     *  to the object returned by this function should not also affect the
     *  NGramMap. This is also known as a "defensive copy". */
    public TimeSeries countHistory(String word) {
        TimeSeries<Double> countingTS = totalMap.get(word);
        return countingTS;
    }

    /** Provides the history of WORD between STARTYEAR and ENDYEAR, inclusive of both ends. The
     *  returned TimeSeries should be a copy, not a link to this NGramMap's TimeSeries. In other words,
     *  changes made to the object returned by this function should not also affect the
     *  NGramMap. This is also known as a "defensive copy". */
    public TimeSeries countHistory(String word, int startYear, int endYear) {
        TimeSeries<Double> countingTS = totalMap.get(word);
        TimeSeries<Double> countingPartsTS = new TimeSeries<>();
        for (int keys: countingTS.keySet()) {
            if (keys >= startYear && keys <= endYear) {
                countingPartsTS.put(keys, countingTS.get(keys));
            }
        }
        return countingPartsTS;
    }

    /** Returns a defensive copy of the total number of words recorded per year in all volumes. */
    public TimeSeries totalCountHistory() {
        TimeSeries<Double> countHistory = new TimeSeries<>();
        countHistory = countsYear;
        return countHistory;
    }

    /** Provides a TimeSeries containing the relative frequency per year of WORD compared to
     *  all words recorded in that year. */
    public TimeSeries weightHistory(String word) {
        /* get the number of the word of that year */
        TimeSeries<Double> wordTS = totalMap.get(word);
        TimeSeries<Double> counts = totalCountHistory();
        TimeSeries<Double> wordWeightTS = new TimeSeries<>();

        for (int key: wordTS.keySet()) {
            wordWeightTS.put(key, wordTS.get(key) / counts.get(key));
        }

        return wordWeightTS;
    }

    /** Provides a TimeSeries containing the relative frequency per year of WORD between STARTYEAR
     *  and ENDYEAR, inclusive of both ends. */
    public TimeSeries weightHistory(String word, int startYear, int endYear) {
        TimeSeries<Double> wordTS = totalMap.get(word);
        TimeSeries<Double> countTS = totalCountHistory();

        TimeSeries<Double> wordTSCut = new TimeSeries<>();
        TimeSeries<Double> countTSCut = new TimeSeries<>();

        for (int year = startYear; year <= endYear; year++) {
            wordTSCut.put(year, wordTS.get(year));
            countTSCut.put(year, countTS.get(year));
        }

        TimeSeries<Double> wordWeight = new TimeSeries<>();

        for (int key: wordTSCut.keySet()) {
            wordWeight.put(key, wordTSCut.get(key) / countTSCut.get(key));
        }

        return wordWeight;
    }

    /** Returns the summed relative frequency per year of all words in WORDS. */
    public TimeSeries summedWeightHistory(Collection<String> words) {
        TimeSeries<Double> wordsSumTS = new TimeSeries<>();
        Iterator wordsIterater = words.iterator();
        wordsSumTS = totalMap.get(wordsIterater);
        wordsIterater.next();

        while (wordsIterater.hasNext()) {
            wordsSumTS.plus(totalMap.get(wordsIterater));
        }

        wordsSumTS.dividedBy(totalCountHistory());

        return wordsSumTS;
    }

    /** Provides the summed relative frequency per year of all words in WORDS
     *  between STARTYEAR and ENDYEAR, inclusive of both ends. If a word does not exist in
     *  this time frame, ignore it rather than throwing an exception. */
    public TimeSeries summedWeightHistory(Collection<String> words,
                                                  int startYear, int endYear) {
        TimeSeries<Double> wordsSumTS = new TimeSeries<>();
        Iterator<String> wordsIterater = words.iterator();
        ArrayList<String> wordsArray = new ArrayList<>();
        while (wordsIterater.hasNext()) {
            wordsArray.add(wordsIterater.next());
        }

        wordsSumTS = totalMap.get(wordsArray.get(0));
        TimeSeries<Double> countTS = totalCountHistory();

        TimeSeries<Double> plused = new TimeSeries<>();
//        TimeSeries<Double> p2 = new TimeSeries<>();
//        TimeSeries<Double> p3 = new TimeSeries<>();
        for (int i = 1; i < wordsArray.size(); i++) {
            plused = wordsSumTS.plus(totalMap.get(wordsArray.get(i)));
            wordsSumTS = plused;
        }

        TimeSeries<Double> wordsSumPart = new TimeSeries<>();
        TimeSeries<Double> countsSumPart = new TimeSeries<>();

        for (int y = startYear; y <= endYear; y++) {
            wordsSumPart.put(y, wordsSumTS.get(y));
            countsSumPart.put(y, countTS.get(y));
        }

        TimeSeries<Double> result = new TimeSeries<>();

        for (int key: wordsSumPart.keySet()) {
            double keyWordSum = wordsSumPart.get(key);
            double keyCountSum = countsSumPart.get(key);
            result.put(key, wordsSumPart.get(key) / countsSumPart.get(key));
        }

        return result;
    }


}
