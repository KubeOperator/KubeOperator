import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MultiClusterBrowserComponent } from './multi-cluster-browser.component';

describe('MultiClusterBrowserComponent', () => {
  let component: MultiClusterBrowserComponent;
  let fixture: ComponentFixture<MultiClusterBrowserComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MultiClusterBrowserComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MultiClusterBrowserComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
