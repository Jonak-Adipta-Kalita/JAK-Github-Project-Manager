import firebase from 'firebase';

const firebaseConfig = {
	apiKey: "AIzaSyAzv968QIYIQjxNdBw0pzoC8BwX1OhiGRw",
	authDomain: "jak--clone-2-7f904.firebaseapp.com",
	projectId: "jak--clone-2-7f904",
	storageBucket: "jak--clone-2-7f904.appspot.com",
	messagingSenderId: "825702221545",
	appId: "1:825702221545:web:d8eb803e06fdfa17f21278"
};

const app = !firebase.apps.length 
	? firebase.initializeApp(firebaseConfig) 
	: firebase.app();
const db = app.firestore()

export default db;
