<!--
  - Copyright (c) 2019 Maximilian Schoenenberg
  -
  - Permission is hereby granted, free of charge, to any person obtaining a copy
  - of this software and associated documentation files (the "Software"), to deal
  - in the Software without restriction, including without limitation the rights
  - to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
  - copies of the Software, and to permit persons to whom the Software is
  - furnished to do so, subject to the following conditions:
  -
  - The above copyright notice and this permission notice shall be included in all
  - copies or substantial portions of the Software.
  -
  - THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
  - IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
  - FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
  - AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
  - LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
  - OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
  - SOFTWARE.
  -->

<template>
    <div id="app" class="container">
        <h1 class="text-center">MedAT-Wortflüssigkeit</h1>
        <div class="my-2 row align-items-end" v-if="myWord">
            <div class="col">
                <word v-bind:shuffled="myWord.Shuffled" v-bind:solution="myWord.Solution">
                    <button class="btn btn-outline-secondary" v-on:click="refresh">Nächstes Wort</button>
                </word>
            </div>
        </div>
        <div class="mt-5 row col justify-content-center">
            <span class=" text-muted">Made with <span style="color: #e25555;">&hearts;</span> for Katja!</span>
        </div>
    </div>
</template>

<script>
    import axios from 'axios'

    import Word from './components/Word';

    export default {
        name: 'app',
        components: {
            Word
        },
        data() {
            return {
                myWord: null
            }
        },
        mounted() {
            this.refresh();
        },
        methods: {
            refresh: function () {
                axios
                    .get("http://localhost:8081/word/new")
                    .then(response => this.myWord = response.data);
            }
        }
    }
</script>
