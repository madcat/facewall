//
//  ViewController.swift
//  facewall-remote
//
//  Created by Lingfei Song on 11/7/15.
//  Copyright © 2015 zealion. All rights reserved.
//

import UIKit

class PrizeViewController: UITableViewController {
    var prizes :[String] = []
    
    @IBOutlet var prizeTableView: UITableView!
    
    override func viewDidLoad() {
        super.viewDidLoad()

        guard let host = NSUserDefaults.standardUserDefaults().stringForKey("host") where host != "" else {
            showError("未设置服务器地址，请设置并重启应用")
            return
        }
        
        let url = NSURL(string: "http://\(host)/prize")
        let request = NSMutableURLRequest(URL: url!, cachePolicy: NSURLRequestCachePolicy.ReloadIgnoringCacheData, timeoutInterval: 20)
        request.HTTPMethod = "GET"
        
        NSURLSession.sharedSession().dataTaskWithRequest(request) { (data, response, error) in
            if error != nil {
                print("GET /prize: \(error?.description)")
                return
            }
            
            guard let httpResponse = response as? NSHTTPURLResponse else {
                self.showError("无法连接服务器，请检查网络并重启应用")
                return
            }
            
            switch httpResponse.statusCode {
            case 200:
                do {
                    if let d = data, let prizes = try NSJSONSerialization.JSONObjectWithData(d, options: NSJSONReadingOptions.MutableContainers) as? [AnyObject] {
                        self.prizes = [String](count:prizes.count, repeatedValue:"")
                        for obj in prizes {
                            if let id = obj["Id"] as? Int, let prize = obj["Prize"] as? String {
                                self.prizes[id-1] = prize
                            }
                        }
                        
                        
                    }
                    dispatch_async(dispatch_get_main_queue(), {
                        self.prizeTableView.reloadData()
                    })
                } catch {
                    self.showError("GET /prize parse json error")
                }
            default:
                self.showError("GET /prize status not 200")
                print("GET /prize HTTP \(httpResponse.statusCode)")
            }
        }.resume()
    }

    override func didReceiveMemoryWarning() {
        super.didReceiveMemoryWarning()
        // Dispose of any resources that can be recreated.
    }
    
    override func numberOfSectionsInTableView(tableView: UITableView) -> Int {
        return 1
    }
    
    override func tableView(tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        return self.prizes.count
    }
    
    override func tableView(tableView: UITableView, cellForRowAtIndexPath indexPath: NSIndexPath) -> UITableViewCell {
        let aCell = self.tableView.dequeueReusableCellWithIdentifier("prizeCell", forIndexPath: indexPath) as! PrizeTableViewCell
        aCell.prizeLabel.text = self.prizes[indexPath.row]
        return aCell
    }
    
    override func tableView(tableView: UITableView, didSelectRowAtIndexPath indexPath: NSIndexPath) {
        self.performSegueWithIdentifier("shuffleForPrize", sender: self)
    }
    
    override func prepareForSegue(segue: UIStoryboardSegue, sender:AnyObject?)
    {
        if (segue.identifier == "shuffleForPrize")
        {
            // upcoming is set to NewViewController (.swift)
            let vc: ShuffleViewController = segue.destinationViewController as! ShuffleViewController
            // indexPath is set to the path that was tapped
            let indexPath = self.prizeTableView.indexPathForSelectedRow!
            // titleString is set to the title at the row in the objects array.
            let prize = self.prizes[indexPath.row]
            // the titleStringViaSegue property of NewViewController is set.
            vc.prize = prize
            vc.title = prize
            self.prizeTableView.deselectRowAtIndexPath(indexPath, animated: true)
        }
    }
    
    func showError(msg :String){
        let alert = UIAlertController(title: "错误", message: msg, preferredStyle: UIAlertControllerStyle.Alert)
        //alert.addAction(UIAlertAction(title: "Click", style: UIAlertActionStyle.Default, handler: nil))
        self.presentViewController(alert, animated: true, completion: nil)
    }
}

