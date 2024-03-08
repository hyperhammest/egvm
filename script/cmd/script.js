var RandomSeed;
var User2Balance = new Map();
var UserList = []
var User2Plants = new Map(); // Plant attribute: height growSpeed lastIrrigateTime
var LastTimestamp = 0;
const BirdEatenHeight = 10;
const FastSeedPrice = 50;
const SlowSeedPrice = 20;
const FastGrownHeightPerHour = 500;
const SlowGrownHeightPerHour = 200;

function setRandomSeed(seedStr) {
	RandomSeed = Sha256(seedStr);
}

function handleInput(inputStr) {
	const input = JSON.parse(inputStr);
	var out = {}
	if(input.op == "depositCoins") {
		out = depositCoins(input);
	} else if(input.op == "buySeed") {
		out = buySeed(input);
	} else if(input.op == "irrigate") {
		out = irrigate(input);
	} else if(input.op == "finishBlock") {
		out = finishBlock(input);
	}
	return JSON.stringify(out);
}

function depositCoins(input) {
	if(User2Balance.has(input.user)) {
		const old = User2Balance.get(input.user);
		User2Balance.set(input.user, input.amount + old);
	} else {
		UserList.push(input.user);
		User2Balance.set(input.user, input.amount);
	}
	return {success: true, input: input}
}

function buySeed(input) {
	if(!User2Balance.has(input.user)) {
		return {success: false, error: "no such user", input: input}
	}
	let which = input.which;
	let price = (FastSeedPrice + SlowSeedPrice) / 2;
	if(which == "fast") {
		price = FastSeedPrice;
	} else if(which == "slow") {
		price = SlowSeedPrice;
	} else { // which == "random"
		const randNum = BufToU32LE(RandomSeed.slice(0,4));
		if(randNum % 1000 < 400) { // pick one randomly
			which = "slow"
		} else {
			which = "fast"
		}
		RandomSeed = Sha256(RandomSeed, Sha256(input.user));
	}
	const growSpeed = which == "slow" ? SlowGrownHeightPerHour : FastGrownHeightPerHour;
	const oldBalance = User2Balance.get(input.user);
	if(oldBalance < price) {
		return {success: false, error: "balance not enough for "+which, input: input}
	}
	User2Balance.set(input.user, oldBalance - price);
	let plants = User2Plants.has(input.user)? User2Plants.get(input.user) : [];
	plants.push({height: 0, growSpeed: growSpeed, lastIrrigateTime: input.timestamp})
	User2Plants.set(input.user, plants);
	return {success: true, extra: {which: which}, input: input}
}

function irrigate(input) {
	if(!User2Plants.has(input.user)) {
		return {success: false, error: "no such user", input: input};
	}
	let plants = User2Plants.get(input.user);
	if(input.plantIndex >= plants.length) {
		return {success: false, error: "no such plant", input: input};
	}
	let plant = plants[input.plantIndex];
	if(input.timestamp < plant.lastIrrigateTime + 3600*24) {
		return {success: false, error: "already irrigated today", input: input};
	}
	plant.lastIrrigateTime = input.timestamp;
	plant.height += plant.growSpeed;
	return {success: true, input: input};
}

function finishBlock(input) {
	RandomSeed = Sha256(RandomSeed, Sha256(input.timestamp.toString()));
	if(input.timestamp < LastTimestamp + 3600 || UserList.length == 0) {
		return {}; //do nothing
	}
	LastTimestamp = input.timestamp;
	const userIndex = BufToU32LE(RandomSeed.slice(0,4)) % UserList.length;
	const user = UserList[userIndex];
	if(!User2Plants.has(user)) {
		return {}; //do nothing
	}
	const plants = User2Plants.get(user);
	if(plants.length == 0) {
		return {}; //do nothing
	}
	const plantIndex = BufToU32LE(RandomSeed.slice(4,8)) % plants.length;
	const plant = plants[plantIndex];
	if(plant.height < BirdEatenHeight) {
		return {}; //do nothing
	}
	plant.height -= BirdEatenHeight;
	return {autoAction: ["bird_eat", {userIndex: userIndex, plantIndex: plantIndex, plant: plant}]}
}

