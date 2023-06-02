export const getTimeDutation=()=>{
    let hour = new Date().getHours()
    if(hour<2){
        return '深夜'
    }else if(hour<6){
        return '凌晨'
    }else if(hour<10){
        return '早上' 
    }else if(hour<12){
        return '上午'
    }else if (hour<14){
        return '中午'
    }else if (hour<18){
        return '下午'
    }else if (hour<22){
        return '晚上'
    }else {
        return '深夜'
    }
}