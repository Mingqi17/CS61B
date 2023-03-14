package ngordnet.ngrams;

import java.util.*;

/** An object for mapping a year number (e.g. 1996) to numerical data. Provides
 *  utility methods useful for data analysis.
 *  @author Josh Hug
 */
public class TimeSeries<L extends Number> extends TreeMap<Integer, Double> {
    /** Constructs a new empty TimeSeries. */

    public TimeSeries() {
        super();
    }

    /** Creates a copy of TS, but only between STARTYEAR and ENDYEAR,
     *  inclusive of both end points. */
    public TimeSeries(TimeSeries<Number> ts, int startYear, int endYear) {
        if (ts != null) {
            for (Integer x: ts.keySet()) {
                if (x >= startYear && x <= endYear) {
                    this.put(x, ts.get(x));
                }
            }
        }
    }

//    public TimeSeries(TimeSeries<T> ts) {
//        for (Integer key : ts.keySet()) {
//            this.put(key, ts.get(key));
//        }
//    }

    /** Returns all years for this TimeSeries (in any order). */
    public List<Integer> years() {
        ArrayList<Integer> years = new ArrayList<>();
        for (Integer x: keySet()) {
            years.add(x);
        }
        return years;
    }

    /** Returns all data for this TimeSeries (in any order).
     *  Must be in the same order as years(). */
    public List<Double> data() {
        ArrayList<Double> vals = new ArrayList<Double>();
        for (Integer x: keySet()) {
            vals.add(get(x));
        }
        return vals;
    }

    /** Returns the yearwise sum of this TimeSeries with the given TS. In other words, for
     *  each year, sum the data from this TimeSeries with the data from TS. Should return a
     *  new TimeSeries (does not modify this TimeSeries). */
    public TimeSeries<Number> plus(TimeSeries<Number> ts) {
        TimeSeries<Number> newSeries = new TimeSeries<Number>();
        ArrayList<Integer> allKeys = new ArrayList<>(this.keySet());
        for (Integer key : ts.keySet()) {
            allKeys.add(key);
        }

        for (Integer x: allKeys) {
            if (ts.get(x) == null) {
                newSeries.put(x, this.get(x));
                continue;
            } else if (this.get(x) == null) {
                newSeries.put(x, ts.get(x));
                continue;
            } else if (this.get(x) == null && ts.get(x) == null) {
                newSeries.put(x, 0.0);
                continue;
            }
            newSeries.put(x, ts.get(x) + this.get(x));
        }

        return newSeries;
    }

     /** Returns the quotient of the value for each year this TimeSeries divided by the
      *  value for the same year in TS. If TS is missing a year that exists in this TimeSeries,
      *  throw an IllegalArgumentException. If TS has a year that is not in this TimeSeries, ignore it.
      *  Should return a new TimeSeries (does not modify this TimeSeries). */
     public TimeSeries<Number> dividedBy(TimeSeries<Number> ts) {
         TimeSeries<Number> newSeries = new TimeSeries<Number>();
         ArrayList<Integer> allKeys = new ArrayList<>(this.keySet());
         for (Integer key : ts.keySet()) {
             allKeys.add(key);
         }

         for (Integer x: allKeys) {
             if (this.containsKey(x) && (!ts.containsKey(x))) {
                 throw new IllegalArgumentException();
             }  else if ((!this.containsKey(x)) && ts.containsKey(x)) {
                 continue;
             }
             newSeries.put(x, this.get(x) / ts.get(x));
         }
         return newSeries;
    }
}
