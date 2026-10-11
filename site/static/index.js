HighlightLisp.highlight_auto();
setTimeout(function(){},0);
var lnheight=12; //start lnheight*XYZ=points lines down
// find with: sed -ne '0,/<pre>/d ; p ; /<\/pre>/Q' index.html | sed -e '/<\|>/d' | wc -l
var lntot=407;
var points=lnheight*lntot;
function rain(){
    if(window.setTimeout&&document.getElementsByTagName){
    var b=document.getElementsByTagName("pre")[0];
    b.style.top=-points+"pt"; //text size
    var l=document.getElementsByTagName("b");
    for(var link in l){
        if(l[link].style)
            l[link].style.top=points+"pt";
    }

    points-=lnheight;

    if(points <= 0)
        points=lnheight*lntot;

    //Speed of use. Lower is faster.
    window.setTimeout(this.rain,2<<6);
    }
}
rain();
