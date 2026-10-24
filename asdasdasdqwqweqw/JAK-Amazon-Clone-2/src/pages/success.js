import Head from "next/head";
import Header from '../components/Header';
import { CheckCircleIcon } from '@heroicons/react/solid';
import { useRouter } from 'next/router';

export default function Success() {
	const router = useRouter();
	return (
		<div className="bg-gray-100 h-screen">
			<Head>
				<title>JAK Amazon Clone 2</title>
				<link rel="icon" href="/favicon.ico" />
			</Head>
			<Header />
			<main className="max-w-screen-lg mx-auto">
				<div className="flex flex-col p-10 bg-white">
					<div className="flex items-center space-x-2 mb-5">
						<CheckCircleIcon className="text-green-500 h-10" />
						<h1 className="text-3xl">Thank You, Your Order has been Confirmed!!</h1>
					</div>
					<p>
						Thank You for Shopping with us. We'll send a Confirmation once your 
						item has shipped. If you would like to check the Status of your Order(s),
						please press the link below.
					</p>
					<button 
						onClick={() => router.push("/orders")} 
						className="button mt-8"
					>
						Go to my Orders
					</button>
				</div>
			</main>
		</div>
	);
};
