import fs from 'node:fs/promises';
import {FileBlob,PresentationFile} from 'file:///C:/Users/iQuantum/.cache/codex-runtimes/codex-primary-runtime/dependencies/node/node_modules/@oai/artifact-tool/dist/artifact_tool.mjs';
const p=await PresentationFile.importPptx(await FileBlob.load('D:/02-A/code/bootcamp/POS_AI_First_Hackathon_2026.pptx'));
const dir='D:/02-A/code/bootcamp/.codex-artifacts/english-demo';
await fs.writeFile(dir+'/source-inspect.ndjson',(await p.inspect({kind:'slide,textbox,image,layout,notes',maxChars:100000})).ndjson);
const rows=[];
for(let i=0;i<p.slides.items.length;i++) {
 const s=p.slides.items[i];
 rows.push({slide:i+1,text:s.shapes.items.map(x=>({id:x.id,text:String(x.text),style:x.text.style,position:x.position}))});
 const b=await p.export({slide:s,format:'png',scale:1});
 await fs.writeFile(dir+`/source-${i+1}.png`,new Uint8Array(await b.arrayBuffer()));
}
await fs.writeFile(dir+'/source-text.json',JSON.stringify(rows,null,2));
console.log('Source rendered',rows.length);
