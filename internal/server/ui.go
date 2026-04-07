package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Quartermaster</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}
body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center;gap:1rem;flex-wrap:wrap}
.hdr h1{font-size:.9rem;letter-spacing:2px}
.hdr h1 span{color:var(--rust)}
.main{padding:1.5rem;max-width:960px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.7rem;text-align:center}
.st-v{font-size:1.3rem;font-weight:700;color:var(--gold)}
.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.2rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;flex-wrap:wrap;align-items:center}
.search{flex:1;min-width:180px;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.filter-sel{padding:.4rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem}
.count-label{font-size:.6rem;color:var(--cm);margin-bottom:.5rem}
.item{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem 1rem;margin-bottom:.5rem;transition:border-color .2s}
.item:hover{border-color:var(--leather)}
.item.low-stock{border-left:3px solid var(--gold)}
.item.out-of-stock{border-left:3px solid var(--red)}
.item-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.8rem}
.item-name{font-size:.85rem;font-weight:700;flex:1}
.item-qty{font-size:1.1rem;font-weight:700;color:var(--cream);white-space:nowrap;text-align:right}
.item-qty-label{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;display:block;margin-top:.1rem}
.item-qty.low{color:var(--gold)}
.item-qty.out{color:var(--red)}
.item-meta{font-size:.55rem;color:var(--cm);margin-top:.4rem;display:flex;gap:.6rem;flex-wrap:wrap;align-items:center}
.item-meta-sep{color:var(--bg3)}
.item-actions{display:flex;gap:.3rem;flex-shrink:0;margin-left:.5rem}
.item-notes{font-size:.65rem;color:var(--cm);margin-top:.4rem;font-style:italic;padding:.3rem .5rem;border-left:2px solid var(--bg3)}
.item-extra{font-size:.58rem;color:var(--cd);margin-top:.4rem;padding-top:.35rem;border-top:1px dashed var(--bg3);display:flex;flex-direction:column;gap:.15rem}
.item-extra-row{display:flex;gap:.4rem}
.item-extra-label{color:var(--cm);text-transform:uppercase;letter-spacing:.5px;min-width:90px}
.item-extra-val{color:var(--cream)}
.badge{font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid var(--bg3);color:var(--cm)}
.badge.new{border-color:var(--blue);color:var(--blue)}
.badge.good{border-color:var(--green);color:var(--green)}
.badge.fair{border-color:var(--gold);color:var(--gold)}
.badge.poor{border-color:var(--red);color:var(--red)}
.badge.broken{border-color:var(--red);color:var(--red)}
.badge.low-stock{border-color:var(--gold);color:var(--gold)}
.badge.out-of-stock{border-color:var(--red);color:var(--red)}
.btn{font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .2s;font-family:var(--mono)}
.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}
.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:480px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}
.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.fr-section{margin-top:1rem;padding-top:.8rem;border-top:1px solid var(--bg3)}
.fr-section-label{font-size:.55rem;color:var(--rust);text-transform:uppercase;letter-spacing:1px;margin-bottom:.5rem}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.row3{display:grid;grid-template-columns:1fr 1fr 1fr;gap:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.85rem}
@media(max-width:600px){.row2,.row3{grid-template-columns:1fr}.toolbar{flex-direction:column;align-items:stretch}.search{min-width:100%}.filter-sel{width:100%}}
</style>
</head>
<body>

<div class="hdr">
<h1 id="dash-title"><span>&#9670;</span> QUARTERMASTER</h1>
<button class="btn btn-p" onclick="openForm()">+ Add Item</button>
</div>

<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar">
<input class="search" id="search" placeholder="Search name, category, location, serial..." oninput="render()">
<select class="filter-sel" id="category-filter" onchange="render()">
<option value="">All Categories</option>
</select>
<select class="filter-sel" id="condition-filter" onchange="render()">
<option value="">All Conditions</option>
<option value="new">New</option>
<option value="good">Good</option>
<option value="fair">Fair</option>
<option value="poor">Poor</option>
<option value="broken">Broken</option>
</select>
<select class="filter-sel" id="location-filter" onchange="render()">
<option value="">All Locations</option>
</select>
</div>
<div class="count-label" id="count"></div>
<div id="list"></div>
</div>

<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()">
<div class="modal" id="mdl"></div>
</div>

<script>
var A='/api';
var RESOURCE='inventory';

// Default categories — overridden by cfg.categories from /api/config if present.
var defaultCategories=['Equipment','Supplies','Tools','Furniture','Electronics','Materials','Other'];

// Field defs drive the form, the rows, and the submit body.
// purchase_price is type 'money' (stored as integer cents).
// quantity is type 'integer'. min_quantity (custom) drives low-stock badges.
var fields=[
{name:'name',label:'Item Name',type:'text',required:true,placeholder:'Cordless drill'},
{name:'quantity',label:'Quantity',type:'integer'},
{name:'category',label:'Category',type:'select',options:defaultCategories},
{name:'location',label:'Location',type:'text',placeholder:'Where it lives'},
{name:'purchase_price',label:'Purchase Price',type:'money'},
{name:'purchase_date',label:'Purchase Date',type:'date'},
{name:'condition',label:'Condition',type:'select',options:['new','good','fair','poor','broken']},
{name:'serial_number',label:'Serial Number',type:'text'},
{name:'notes',label:'Notes',type:'textarea'}
];

var items=[],editId=null;

// ─── Helpers ──────────────────────────────────────────────────────

function fmtMoney(cents){
if(cents===null||cents===undefined||cents==='')return'$0.00';
var n=parseInt(cents,10);
if(isNaN(n))return'$0.00';
var neg=n<0;
n=Math.abs(n);
var dollars=Math.floor(n/100);
var rem=n%100;
var s='$'+dollars.toLocaleString()+'.'+(rem<10?'0':'')+rem;
return neg?'-'+s:s;
}

function parseMoney(str){
if(!str)return 0;
var s=String(str).replace(/[^0-9.\-]/g,'');
var f=parseFloat(s);
if(isNaN(f))return 0;
return Math.round(f*100);
}

function fmtDate(s){
if(!s)return'';
try{return new Date(s).toLocaleDateString('en-US',{month:'short',day:'numeric',year:'numeric'})}catch(e){return s}
}

// Determine stock state for an item using its min_quantity custom field (if set).
// Returns 'out' (qty=0), 'low' (qty<min and min>0), or '' (normal).
function stockState(item){
var qty=parseInt(item.quantity||0,10);
if(qty===0)return'out';
var min=parseInt(item.min_quantity||0,10);
if(min>0&&qty<=min)return'low';
return'';
}

// ─── Loading and rendering ────────────────────────────────────────

async function load(){
try{
var r=await fetch(A+'/'+RESOURCE).then(function(r){return r.json()});
var list=r[RESOURCE]||[];
try{
var extras=await fetch(A+'/extras/'+RESOURCE).then(function(r){return r.json()});
list.forEach(function(it){
var ex=extras[it.id];
if(!ex)return;
Object.keys(ex).forEach(function(k){if(it[k]===undefined)it[k]=ex[k]});
});
}catch(e){}
items=list;
}catch(e){
console.error('load failed',e);
items=[];
}
populateFilters();
renderStats();
render();
}

function populateFilters(){
var catSel=document.getElementById('category-filter');
var locSel=document.getElementById('location-filter');
var seenCat={},seenLoc={};
var cats=[],locs=[];
var catField=fieldByName('category');
if(catField&&catField.options){catField.options.forEach(function(c){if(!seenCat[c]){seenCat[c]=true;cats.push(c)}})}
items.forEach(function(i){
if(i.category&&!seenCat[i.category]){seenCat[i.category]=true;cats.push(i.category)}
if(i.location&&!seenLoc[i.location]){seenLoc[i.location]=true;locs.push(i.location)}
});
if(catSel){
var cur=catSel.value;
catSel.innerHTML='<option value="">All Categories</option>'+cats.map(function(c){return'<option value="'+esc(c)+'"'+(c===cur?' selected':'')+'>'+esc(c)+'</option>'}).join('');
}
if(locSel){
var cur2=locSel.value;
locSel.innerHTML='<option value="">All Locations</option>'+locs.map(function(l){return'<option value="'+esc(l)+'"'+(l===cur2?' selected':'')+'>'+esc(l)+'</option>'}).join('');
}
}

function renderStats(){
var total=items.length;
var totalQty=0;
var totalValue=0;
items.forEach(function(i){
var q=parseInt(i.quantity||0,10);
var p=parseInt(i.purchase_price||0,10);
totalQty+=q;
totalValue+=q*p;
});
document.getElementById('stats').innerHTML=
'<div class="st"><div class="st-v">'+total+'</div><div class="st-l">Items</div></div>'+
'<div class="st"><div class="st-v">'+totalQty.toLocaleString()+'</div><div class="st-l">Total Qty</div></div>'+
'<div class="st"><div class="st-v">'+fmtMoney(totalValue)+'</div><div class="st-l">Total Value</div></div>';
}

function render(){
var q=(document.getElementById('search').value||'').toLowerCase();
var cf=document.getElementById('category-filter').value;
var conf=document.getElementById('condition-filter').value;
var lf=document.getElementById('location-filter').value;
var f=items;
if(cf)f=f.filter(function(i){return i.category===cf});
if(conf)f=f.filter(function(i){return i.condition===conf});
if(lf)f=f.filter(function(i){return i.location===lf});
if(q)f=f.filter(function(i){
return(i.name||'').toLowerCase().includes(q)||
       (i.category||'').toLowerCase().includes(q)||
       (i.location||'').toLowerCase().includes(q)||
       (i.serial_number||'').toLowerCase().includes(q);
});
document.getElementById('count').textContent=f.length+' item'+(f.length!==1?'s':'');
if(!f.length){
var msg=window._emptyMsg||'No items found.';
document.getElementById('list').innerHTML='<div class="empty">'+esc(msg)+'</div>';
return;
}
var h='';
f.forEach(function(i){h+=itemHTML(i)});
document.getElementById('list').innerHTML=h;
}

function itemHTML(i){
var stock=stockState(i);
var cls='item';
if(stock==='out')cls+=' out-of-stock';
else if(stock==='low')cls+=' low-stock';

var qtyCls='item-qty';
if(stock==='out')qtyCls+=' out';
else if(stock==='low')qtyCls+=' low';

var h='<div class="'+cls+'"><div class="item-top">';
h+='<div class="item-name">'+esc(i.name)+'</div>';
h+='<div class="'+qtyCls+'">'+esc(String(i.quantity||0))+'<span class="item-qty-label">in stock</span></div>';
h+='<div class="item-actions">';
h+='<button class="btn btn-sm" onclick="openEdit(\''+i.id+'\')">Edit</button>';
h+='<button class="btn btn-sm" onclick="del(\''+i.id+'\')" style="color:var(--red)">&#10005;</button>';
h+='</div></div>';

h+='<div class="item-meta">';
var parts=[];
if(i.category)parts.push('<span>'+esc(i.category)+'</span>');
if(i.location)parts.push('<span>'+esc(i.location)+'</span>');
if(i.purchase_price)parts.push('<span>'+esc(fmtMoney(i.purchase_price))+'</span>');
if(i.purchase_date)parts.push('<span>'+esc(fmtDate(i.purchase_date))+'</span>');
if(i.serial_number)parts.push('<span>SN: '+esc(i.serial_number)+'</span>');
h+=parts.join('<span class="item-meta-sep">·</span>');
if(i.condition)h+=' <span class="badge '+esc(i.condition)+'">'+esc(i.condition)+'</span>';
if(stock==='out')h+=' <span class="badge out-of-stock">OUT OF STOCK</span>';
else if(stock==='low')h+=' <span class="badge low-stock">LOW STOCK</span>';
h+='</div>';

if(i.notes)h+='<div class="item-notes">'+esc(i.notes)+'</div>';

// Custom fields from personalization (excluding min_quantity which we already
// use for the low-stock state — but we still want to show it)
var customRows='';
fields.forEach(function(f){
if(!f.isCustom)return;
var v=i[f.name];
if(v===undefined||v===null||v==='')return;
customRows+='<div class="item-extra-row">';
customRows+='<span class="item-extra-label">'+esc(f.label)+'</span>';
customRows+='<span class="item-extra-val">'+esc(String(v))+'</span>';
customRows+='</div>';
});
if(customRows)h+='<div class="item-extra">'+customRows+'</div>';

h+='</div>';
return h;
}

// ─── Form ─────────────────────────────────────────────────────────

function fieldByName(n){
for(var i=0;i<fields.length;i++)if(fields[i].name===n)return fields[i];
return null;
}

function fieldHTML(f,value){
var v=value;
if(v===undefined||v===null)v='';
var req=f.required?' *':'';
var ph='';
if(f.placeholder)ph=' placeholder="'+esc(f.placeholder)+'"';
else if(f.name==='name'&&window._placeholderName)ph=' placeholder="'+esc(window._placeholderName)+'"';

var h='<div class="fr"><label>'+esc(f.label)+req+'</label>';

if(f.type==='select'){
h+='<select id="f-'+f.name+'">';
if(!f.required)h+='<option value="">Select...</option>';
(f.options||[]).forEach(function(o){
var sel=(String(v)===String(o))?' selected':'';
var disp=(typeof o==='string')?(o.charAt(0).toUpperCase()+o.slice(1)):String(o);
h+='<option value="'+esc(String(o))+'"'+sel+'>'+esc(disp)+'</option>';
});
h+='</select>';
}else if(f.type==='textarea'){
h+='<textarea id="f-'+f.name+'" rows="2"'+ph+'>'+esc(String(v))+'</textarea>';
}else if(f.type==='checkbox'){
h+='<input type="checkbox" id="f-'+f.name+'"'+(v?' checked':'')+' style="width:auto">';
}else if(f.type==='money'){
var displayVal=v?(parseInt(v,10)/100).toFixed(2):'';
h+='<input type="text" id="f-'+f.name+'" value="'+esc(displayVal)+'"'+ph+' inputmode="decimal">';
}else if(f.type==='integer'||f.type==='number'){
h+='<input type="number" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
}else{
var inputType=f.type||'text';
h+='<input type="'+esc(inputType)+'" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
}

h+='</div>';
return h;
}

function formHTML(item){
var i=item||{};
var isEdit=!!item;
var h='<h2>'+(isEdit?'EDIT ITEM':'NEW ITEM')+'</h2>';

// Name on its own row
h+=fieldHTML(fieldByName('name'),i.name);

// Quantity + category on one row
h+='<div class="row2">'+fieldHTML(fieldByName('quantity'),i.quantity)+fieldHTML(fieldByName('category'),i.category)+'</div>';

// Location + condition
h+='<div class="row2">'+fieldHTML(fieldByName('location'),i.location)+fieldHTML(fieldByName('condition'),i.condition)+'</div>';

// Purchase price + purchase date
h+='<div class="row2">'+fieldHTML(fieldByName('purchase_price'),i.purchase_price)+fieldHTML(fieldByName('purchase_date'),i.purchase_date)+'</div>';

// Serial number on its own (often long)
h+=fieldHTML(fieldByName('serial_number'),i.serial_number);

// Notes
h+=fieldHTML(fieldByName('notes'),i.notes);

// Custom fields injected by personalization
var customFields=fields.filter(function(f){return f.isCustom});
if(customFields.length){
var sectionLabel=window._customSectionLabel||'Additional Details';
h+='<div class="fr-section"><div class="fr-section-label">'+esc(sectionLabel)+'</div>';
customFields.forEach(function(f){h+=fieldHTML(f,i[f.name])});
h+='</div>';
}

h+='<div class="acts">';
h+='<button class="btn" onclick="closeModal()">Cancel</button>';
h+='<button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Add Item')+'</button>';
h+='</div>';
return h;
}

function openForm(){
editId=null;
document.getElementById('mdl').innerHTML=formHTML();
document.getElementById('mbg').classList.add('open');
var n=document.getElementById('f-name');
if(n)n.focus();
}

function openEdit(id){
var x=null;
for(var j=0;j<items.length;j++){if(items[j].id===id){x=items[j];break}}
if(!x)return;
editId=id;
document.getElementById('mdl').innerHTML=formHTML(x);
document.getElementById('mbg').classList.add('open');
}

function closeModal(){
document.getElementById('mbg').classList.remove('open');
editId=null;
}

// ─── Submit ───────────────────────────────────────────────────────

async function submit(){
var nameEl=document.getElementById('f-name');
if(!nameEl||!nameEl.value.trim()){alert('Item name is required');return}

var body={};
var extras={};
fields.forEach(function(f){
var el=document.getElementById('f-'+f.name);
if(!el)return;
var val;
if(f.type==='checkbox')val=el.checked;
else if(f.type==='money')val=parseMoney(el.value);
else if(f.type==='integer')val=parseInt(el.value,10)||0;
else if(f.type==='number')val=parseFloat(el.value)||0;
else val=el.value.trim();
if(f.isCustom)extras[f.name]=val;
else body[f.name]=val;
});

var savedId=editId;
try{
if(editId){
var r1=await fetch(A+'/'+RESOURCE+'/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r1.ok){var e1=await r1.json().catch(function(){return{}});alert(e1.error||'Save failed');return}
}else{
var r2=await fetch(A+'/'+RESOURCE,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r2.ok){var e2=await r2.json().catch(function(){return{}});alert(e2.error||'Save failed');return}
var created=await r2.json();
savedId=created.id;
}
if(savedId&&Object.keys(extras).length){
await fetch(A+'/extras/'+RESOURCE+'/'+savedId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(extras)}).catch(function(){});
}
}catch(e){
alert('Network error: '+e.message);
return;
}

closeModal();
load();
}

async function del(id){
if(!confirm('Delete this item?'))return;
await fetch(A+'/'+RESOURCE+'/'+id,{method:'DELETE'});
load();
}

function esc(s){
if(s===undefined||s===null)return'';
var d=document.createElement('div');
d.textContent=String(s);
return d.innerHTML;
}

document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal()});

// ─── Personalization ──────────────────────────────────────────────

(function loadPersonalization(){
fetch('/api/config').then(function(r){return r.json()}).then(function(cfg){
if(!cfg||typeof cfg!=='object')return;

if(cfg.dashboard_title){
var h1=document.getElementById('dash-title');
if(h1)h1.innerHTML='<span>&#9670;</span> '+esc(cfg.dashboard_title);
document.title=cfg.dashboard_title;
}

if(cfg.empty_state_message)window._emptyMsg=cfg.empty_state_message;
if(cfg.placeholder_name)window._placeholderName=cfg.placeholder_name;
if(cfg.primary_label)window._customSectionLabel=cfg.primary_label+' Details';

// Categories from config replace the default category options
if(Array.isArray(cfg.categories)&&cfg.categories.length){
var catField=fieldByName('category');
if(catField)catField.options=cfg.categories.slice();
}

if(Array.isArray(cfg.custom_fields)){
cfg.custom_fields.forEach(function(cf){
if(!cf||!cf.name||!cf.label)return;
if(fieldByName(cf.name))return;
fields.push({
name:cf.name,
label:cf.label,
type:cf.type||'text',
options:cf.options||[],
isCustom:true
});
});
}
}).catch(function(){
}).finally(function(){
load();
});
})();
</script>
</body>
</html>`
