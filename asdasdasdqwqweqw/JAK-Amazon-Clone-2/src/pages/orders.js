import Head from "next/head";
import Header from '../components/Header';
import { getSession, useSession } from 'next-auth/client';
import db from '../../firebase';
import moment from 'moment';
import Order from '../components/Order';

export default function Orders({ orders }) {
	const [session] = useSession();
	return (
		<div className="">
			<Head>
				<title>JAK Amazon Clone 2</title>
				<link rel="icon" href="/favicon.ico" />
			</Head>
			<Header />
			<main className="max-w-screen-lg mx-auto p-10">
				{session ? (
					<>
						<h1 className="text-3xl border-b mb-2 pb-1 border-yellow-400">
							Your Orders
						</h1>
						<h2>{orders.length} Order(s)</h2>
					</>
				): (
					<h2 className="text-center text-2xl">
						Please Sign In to see your Orders!!
					</h2>
				)}
				<div className="mt-5 space-y-4">
					{orders?.map(({ id, amount, amountShipping, items, timestamp, images }) => (
						<Order 
							key={id}
							id={id} 
							amount={amount} 
							amountShipping={amountShipping} 
							items={items} 
							timestamp={timestamp} 
							images={images}
						/>
					))}
				</div>
			</main>
		</div>
	);
};

export async function getServerSideProps(context) {
	const stripe = require("stripe")(process.env.STRIPE_SECRET_KEY);
	const session = await getSession(context);
	if (!session) {
		return {
			props: {},
		};
	};
	const stripeOrders = await db
		.collection('users')
		.doc(session.user.email)
		.collection("orders")
		.orderBy('timestamp', 'desc')
		.get();
	const orders = await Promise.all(
		stripeOrders.docs.map(async (order) => ({
			id: order.id,
			amount: order.data().amount,
			amountShipping: order.data().amount_shipping,
			images: order.data().images,
			timestamp: moment(order.data().timestamp.toDate()).unix(),
			items: (
				await stripe.checkout.sessions.listLineItems(order.id, {
					limit: 100,
				})
			).data,
		}))
	);
	return {
		props: {
			orders,
		},
	};
};
