package ngordnet.main;

import ngordnet.hugbrowsermagic.NgordnetQuery;
import ngordnet.hugbrowsermagic.NgordnetQueryHandler;
import ngordnet.ngrams.NGramMap;
import ngordnet.ngrams.TimeSeries;
import ngordnet.plotting.Plotter;
import org.knowm.xchart.XYChart;

import java.sql.Time;
import java.util.ArrayList;


public class HistoryHandler extends NgordnetQueryHandler {

    NGramMap mapMap;
    public HistoryHandler(NGramMap map) {
        mapMap = map;
    }

    public String handle(NgordnetQuery q) {

//        TimeSeries curWordData = new TimeSeries();

        ArrayList<TimeSeries<Number>> lts = new ArrayList<>();
        ArrayList<String> labels = new ArrayList<>();

        for (int i = 0; i < q.words().size(); i++) {
            labels.add(q.words().get(i));
            lts.add(mapMap.countHistory(q.words().get(i), q.startYear(), q.endYear()));
        }

        XYChart chart = Plotter.generateTimeSeriesChart(labels, lts);
        String encodedImage = Plotter.encodeChartAsString(chart);

        return encodedImage;
    }
}
