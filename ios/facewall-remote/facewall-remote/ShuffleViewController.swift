//
//  ShuffleViewController.swift
//  facewall-remote
//
//  Created by Lingfei Song on 11/7/15.
//  Copyright © 2015 zealion. All rights reserved.
//

import UIKit

struct Winner {
    var Step:Int
    var Prize:String
    var Code:String
    var Tag:String
}

class ShuffleViewController: UIViewController {
    
    var prize:String = ""
    var winners:[Winner] = []
    
    @IBOutlet weak var shuffleButton: UIButton!
    @IBOutlet weak var stepsTextView: UITextView!
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        self.shuffleButton.layer.cornerRadius = self.shuffleButton.frame.width/2
        self.updateResult()
    }

    override func didReceiveMemoryWarning() {
        super.didReceiveMemoryWarning()
        // Dispose of any resources that can be recreated.
    }
    
    func updateResult(){
        self.shuffleButton.hidden = true
        guard let host = NSUserDefaults.standardUserDefaults().stringForKey("host") where host != "" else {
            showError("未设置服务器地址，请设置并重启应用")
            return
        }
        
        let url = NSURL(string: "http://\(host)/winner")
        let request = NSMutableURLRequest(URL: url!, cachePolicy: NSURLRequestCachePolicy.ReloadIgnoringCacheData, timeoutInterval: 20)
        request.HTTPMethod = "GET"
        
        NSURLSession.sharedSession().dataTaskWithRequest(request) { (data, response, error) in
            if error != nil {
                print("GET /winner: \(error?.description)")
                return
            }
            
            guard let httpResponse = response as? NSHTTPURLResponse else {
                self.showError("无法连接服务器，请检查网络并重启应用")
                return
            }
            
            switch httpResponse.statusCode {
            case 200:
                do {
                    if let d = data, let winners = try NSJSONSerialization.JSONObjectWithData(d, options: NSJSONReadingOptions.MutableContainers) as? [AnyObject] {
                        
                        for obj in winners {
                            if let step = obj["Step"] as? Int, let prize = obj["Prize"] as? String, let code = obj["Code"] as? String, let tag = obj["Tag"] as? String {
                                let w = Winner(Step: step, Prize: prize, Code: code, Tag: tag)
                                self.winners.append(w)
                            }
                        }
                        
                        var steps = [Int:Array<Winner>]()
                        for w in self.winners {
                            if w.Prize != self.prize {continue}
                            if steps[w.Step] == nil {
                                steps[w.Step] = [Winner]()
                            }
                            steps[w.Step]!.append(w)
                        }
                        
                        for (k,v) in Array(steps).sort({$0.0 > $1.0}) {
                            dispatch_async(dispatch_get_main_queue(), {
                                self.stepsTextView.text = self.stepsTextView.text + "\n== 第\(k)次: \(self.prize) ==\n"
                                for w in v as [Winner] {
                                    self.stepsTextView.text = self.stepsTextView.text + "\(w.Code) (\(w.Tag))  "
                                }
                                self.stepsTextView.text = self.stepsTextView.text + "\n"
                            })
                        }
                        dispatch_async(dispatch_get_main_queue(), {
                            self.shuffleButton.hidden = false
                        })
                    }
                } catch {
                    self.showError("GET /winner parse json error")
                }
            default:
                self.showError("GET /prize status not 200")
                print("GET /winner HTTP \(httpResponse.statusCode)")
            }
            }.resume()
    }
    

    /*
    // MARK: - Navigation

    // In a storyboard-based application, you will often want to do a little preparation before navigation
    override func prepareForSegue(segue: UIStoryboardSegue, sender: AnyObject?) {
        // Get the new view controller using segue.destinationViewController.
        // Pass the selected object to the new view controller.
    }
    */
    
    func showError(msg :String){
        let alert = UIAlertController(title: "错误", message: msg, preferredStyle: UIAlertControllerStyle.Alert)
        //alert.addAction(UIAlertAction(title: "Click", style: UIAlertActionStyle.Default, handler: nil))
        self.presentViewController(alert, animated: true, completion: nil)
    }

}
