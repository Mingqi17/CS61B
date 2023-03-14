package ngordnet.main;

import ngordnet.ngrams.NGramMap;
import ngordnet.hugbrowsermagic.HugNgordnetServer;

public class Main {
    public static void main(String[] args) {

        HugNgordnetServer hns = new HugNgordnetServer();

        hns.startUp();

        String wordFile = "./data/ngrams/top_14377_words.csv";
        String countFile = "./data/ngrams/total_counts.csv";
        NGramMap ngm = new NGramMap(wordFile, countFile);

        /* NgordnetQuery is what being typed, which are words, startYear, endYear */
        hns.register("historytext", new HistoryTextHandler(ngm));
        hns.register("history", new HistoryHandler(ngm));
    }
}
