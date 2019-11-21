clear; close all;

filename = "../results/50x50.json";

res = jsondecode(fileread(filename));

path = res.Path + 1;
width = res.MG.Terrain.Width;
height = res.MG.Terrain.Height;

drones = res.MG.Drones;
visibilitySlices = res.MG.VisibilitySlices;

startTime = 0;
endTime = drones(1).EndTime;

% for droneID = 1:length(drones)
%     if drones(droneID).EndTime > endTime
%         endTime = drones(droneID).EndTime;
%     end
% end

movie = [];

maxTerrHeight = max(max(res.MG.Terrain.Heights));
c = 0:0.01:1;
c = c .* maxTerrHeight;
r1 = 13;
g1 = 115;
b1 = 0;
r2 = 69;
g2 = 40;
b2 = 0;
C = [r1 + (r2 - r1) / maxTerrHeight * c; ...
     g1 + (g2 - g1) / maxTerrHeight * c; ...
     b1 + (b2 - b1) / maxTerrHeight * c]' / 225;
% colormap(C); % set the colors to green

for time = 0:length(path)-1
    visibilityAtTime = visibilitySlices(time+1).Visible;
    h = figure;
    hold on;
    G = res.MG.Terrain.Heights;
%     C = colormap;  % Get the figure's colormap.
    L = size(C,1);
    Gs = round(interp1(linspace(min(G(:)),max(G(:)),L),1:L,G));
    H = reshape(C(Gs,:),[size(Gs) 3]);
%     rgbImage = ind2gray(res.MG.Terrain.Heights, hsv(100));
    
    visibilityRegion = H;
    [rows, cols, vs] = find(visibilityAtTime);
    
    for i=1:length(rows)
        visibilityRegion(rows(i), cols(i), :) = [1 0 0];
    end
    
    H = visibilityRegion;
    
    s = surf(G, H);
%     colormap("summer")
    plot3(path(time+1, 1), path(time+1, 2), G(path(time+1, 1), path(time+1, 2)) + 1, '.c', 'MarkerSize', 30);
    title(['t = ' num2str(time+1)])
    
    for droneID = 1:length(drones)
        drone = drones(droneID);
        velocity = [drone.Velocity.X, drone.Velocity.Y];
        startPosition = [drone.Origin.X, drone.Origin.Y];
        endPosition = [drone.Destination.X, drone.Destination.Y];
        if time >= drone.StartTime && time <= drone.EndTime
            position = (time-drone.StartTime) * velocity + startPosition + 1;
            plot3(position(1), position(2), 1, '.y', 'MarkerSize', 30);
        end
    end
    
    view([30 30]);
    movie = [movie, getframe(h)];
    zlim([0 35]);
end

