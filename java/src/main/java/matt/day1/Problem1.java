package matt.day1;

import java.io.File;
import java.io.IOException;
import java.net.URISyntaxException;
import java.nio.file.Files;
import java.util.stream.Stream;

public class Problem1 {

    private static int readInput(File f) {
        int freq = 0;
        try (Stream<String> lines = Files.lines(f.toPath())) {
            freq = lines.mapToInt(n -> Integer.parseInt(n)).sum();
        } catch (IOException e) {
            e.printStackTrace();
        }
        return freq;
    }

    private static void operate(int i) {

    }
    public static void main(String args[]) {
        try {
            int freq = readInput(new File(Problem1.class.getResource("/day1").toURI()));
            System.out.println(freq);
        } catch (URISyntaxException e) {
            e.printStackTrace();
        }
    }
}
