// copied from CMU Sphinx Samples
// typical invocation:
//
// javac Transcriber.java -cp /sphinx4-5prealpha-src/sphinx4-core/build/libs/sphinx4-core-5prealpha-SNAPSHOT.jar
// java -Xmx2g -cp /sphinx4-5prealpha-src/sphinx4-core/build/libs/sphinx4-core-5prealpha-SNAPSHOT.jar:/sphinx4-5prealpha-src/sphinx4-data/build/libs/sphinx4-data-5prealpha-SNAPSHOT.jar:. Transcriber meeting.wav transcription.txt
//
// We don't need no package 'cause we're bad boys
// package com.example;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.InputStream;
import java.io.OutputStream;
import java.nio.charset.StandardCharsets;

import edu.cmu.sphinx.api.Configuration;
import edu.cmu.sphinx.api.SpeechResult;
import edu.cmu.sphinx.api.StreamSpeechRecognizer;

public class Transcriber {

    public static void main(String[] args) throws Exception {

        String infile = args[0];
        String outfile = args[1];

        Configuration configuration = new Configuration();

        configuration
                .setAcousticModelPath("resource:/edu/cmu/sphinx/models/en-us/en-us");
        configuration
                .setDictionaryPath("resource:/edu/cmu/sphinx/models/en-us/cmudict-en-us.dict");
        configuration
                .setLanguageModelPath("resource:/edu/cmu/sphinx/models/en-us/en-us.lm.bin");

        StreamSpeechRecognizer recognizer = new StreamSpeechRecognizer(
                configuration);
        InputStream input  = new FileInputStream(new File(infile));
        OutputStream output = new FileOutputStream(new File(outfile));

        recognizer.startRecognition(input);
        SpeechResult result;
        while ((result = recognizer.getResult()) != null) {
            output.write(result.getHypothesis().getBytes(StandardCharsets.UTF_8));
        }
        recognizer.stopRecognition();
    }
}
