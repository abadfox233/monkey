fn(x){
	if(x<3){
		return 1;
	}
	prev = 1;
	curr = 1;
	for(i = 2;i<x;i = i+1;){
		temp = curr;
		curr = prev + curr;
		prev = temp;
	}
	return curr;
}